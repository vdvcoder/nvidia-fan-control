package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

// Config structure to match the provided JSON
type Config struct {
	TimeToUpdate      int                `json:"time_to_update"`
	TemperatureRanges []TemperatureRange `json:"temperature_ranges"`
}

type TemperatureRange struct {
	MinTemperature int `json:"min_temperature"`
	MaxTemperature int `json:"max_temperature"`
	FanSpeed       int `json:"fan_speed"`
	Hysteresis     int `json:"hysteresis"`
}

// LoadConfig reads the JSON config file
func loadConfig(file string) (Config, error) {
	var config Config
	data, err := os.ReadFile(file) // Updated to use os.ReadFile
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

// Absolute difference function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetFanSpeedForTemperature determines the appropriate fan speed based on temperature and hysteresis
func getFanSpeedForTemperature(temp, prevTemp, prevSpeed int, ranges []TemperatureRange) int {
	for _, r := range ranges {
		if temp > r.MinTemperature && temp <= r.MaxTemperature {
			// Apply hysteresis: Change fan speed only if temperature has moved significantly
			if abs(temp-prevTemp) >= r.Hysteresis || prevSpeed != r.FanSpeed {
				return r.FanSpeed
			}
		}
	}
	// Default to the previous fan speed if no range matches
	return prevSpeed
}

func main() {
	// Open a log file
	logFile, err := os.OpenFile("/var/log/nvidia_fan_control.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set logging to file
	log.SetOutput(logFile)

	// Load configuration
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize NVML
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()

	// Get GPU count
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", nvml.ErrorString(ret))
	}

	// Initialize tracking variables
	prevTemps := make([]int, count)     // Previous temperatures for each GPU
	prevFanSpeeds := make([]int, count) // Previous fan speeds for each GPU

	// Monitoring loop
	for {
		for i := 0; i < count; i++ {
			device, ret := nvml.DeviceGetHandleByIndex(i)
			if ret != nvml.SUCCESS {
				log.Printf("Unable to get device at index %d: %v", i, nvml.ErrorString(ret))
				continue
			}

			// Get current temperature
			temp, ret := nvml.DeviceGetTemperature(device, nvml.TEMPERATURE_GPU)
			if ret != nvml.SUCCESS {
				log.Printf("Unable to get temperature for GPU %d: %v", i, nvml.ErrorString(ret))
				continue
			}

			// Convert uint32 temperature to int for compatibility
			tempInt := int(temp)

			// Determine appropriate fan speed
			newFanSpeed := getFanSpeedForTemperature(tempInt, prevTemps[i], prevFanSpeeds[i], config.TemperatureRanges)

			// Update fan speed if it has changed
			if newFanSpeed != prevFanSpeeds[i] {
				// Set manual fan control policy
				ret = nvml.DeviceSetFanControlPolicy(device, 0, 1)
				if ret != nvml.SUCCESS {
					log.Printf("Unable to set manual fan control policy for GPU %d: %v", i, nvml.ErrorString(ret))
					continue
				}

				// Set the new fan speed
				ret = nvml.DeviceSetFanSpeed_v2(device, 0, newFanSpeed)
				if ret != nvml.SUCCESS {
					log.Printf("Unable to set fan speed for GPU %d: %v", i, nvml.ErrorString(ret))
					continue
				}

				log.Printf("Updated GPU %d: Temp=%d°C, Fan Speed=%d%%", i, tempInt, newFanSpeed)
				prevFanSpeeds[i] = newFanSpeed
			} else {
				log.Printf("GPU %d: Temp=%d°C, Fan Speed=%d%% (unchanged)", i, tempInt, newFanSpeed)
			}

			// Update the previous temperature
			prevTemps[i] = tempInt
		}

		// Wait for the configured update interval
		time.Sleep(time.Duration(config.TimeToUpdate) * time.Second)
	}
}
