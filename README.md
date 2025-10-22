# Nvidia Fan Control

# Setup for Omarchy 3.1

A lightweight Linux utility for monitoring GPU temperatures and dynamically controlling NVIDIA GPU fan speeds using NVML.

## Requirements
- NVIDIA GPUs with NVML support
- NVIDIA drivers 520 or higher

## Download
```bash
cd ~/.config

git clone git@github.com:vdvcoder/nvidia-fan-control.git

```

## Build 
```bash
cd ~/.config

cd nvidia-fan-control

go build -o nvidia-fan-control
```

## Installation
```bash
sudo mv nvidia-fan-control /usr/local/bin/
```

## Configuration
edit the file `config.json` with the following structure
```
{
    "time_to_update": 5,
    "temperature_ranges": [
      { "min_temperature": 0, "max_temperature": 40, "fan_speed": 30, "hysteresis": 3 },
      { "min_temperature": 40, "max_temperature": 60, "fan_speed": 40, "hysteresis": 3 },
      { "min_temperature": 60, "max_temperature": 80, "fan_speed": 70, "hysteresis": 3 },
      { "min_temperature": 80, "max_temperature": 100, "fan_speed": 100, "hysteresis": 3 },
      { "min_temperature": 100, "max_temperature": 200, "fan_speed": 100, "hysteresis": 0 }
    ]
  }
```

## Service
```bash
sudo nano /etc/systemd/system/nvidia-fan-control.service
```
update WorkingDirectory and set the path to your config file
```
[Unit]
Description=NVIDIA Fan Control Service
After=network.target

[Service]
ExecStart=/usr/local/bin/nvidia-fan-control
WorkingDirectory=~/config/nvidia-fan-control
StandardOutput=file:/var/log/nvidia_fan_control.log
StandardError=file:/var/log/nvidia_fan_control_error.log
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable nvidia-fan-control.service
sudo systemctl start nvidia-fan-control.service
sudo systemctl status nvidia-fan-control.service
```

### Check Logs
```bash
sudo tail -f /var/log/nvidia_fan_control.log
```
