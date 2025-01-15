# Nvidia Fan Control
Nvidia Fan Control for linux

## Build
```bash
go build -o nvidia_fan_control
```

## Installation
```bash
sudo mv nvidia_fan_control /usr/local/bin/
```

## Configuration
edit the file `config.json` with the following structure
```
{
    "time_to_update": 5,
    "temperature_ranges": [
      { "min_temperature": 0, "max_temperature": 40, "fan_speed": 30, "hysteresis": 2 },
      { "min_temperature": 40, "max_temperature": 60, "fan_speed": 40, "hysteresis": 2 },
      { "min_temperature": 60, "max_temperature": 80, "fan_speed": 70, "hysteresis": 2 },
      { "min_temperature": 80, "max_temperature": 100, "fan_speed": 100, "hysteresis": 2 },
      { "min_temperature": 100, "max_temperature": 200, "fan_speed": 100, "hysteresis": 0 }
    ]
  }
```

## Service
```bash
/etc/systemd/system/nvidia_fan_control.service
```

```
[Unit]
Description=NVIDIA Fan Control Service
After=network.target

[Service]
ExecStart=/usr/bin/sudo /path/to/your/executable
WorkingDirectory=/path/to/your/config
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
sudo systemctl enable nvidia_fan_control.service
sudo systemctl start nvidia_fan_control.service
sudo systemctl status nvidia_fan_control.service
```
### Check Logs
```bash
sudo tail -f /var/log/nvidia_fan_control.log
```

## Drivers
It works with nvidia drivers 520 or higher