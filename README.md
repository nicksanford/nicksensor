# Boom Sensor

### Build:
```
make
```

### Example Config:
```json
{
  "components": [
    {
      "name": "boom",
      "namespace": "rdk",
      "type": "sensor",
      "model": "ncs:sensor:boom",
      "attributes": {},
      "service_configs": [
        {
          "type": "data_manager",
          "attributes": {
            "capture_methods": [
              {
                "method": "Readings",
                "capture_frequency_hz": 1,
                "additional_params": {}
              }
            ]
          }
        }
      ]
    }
  ],
  "services": [
    {
      "name": "data_manager-1",
      "namespace": "rdk",
      "type": "data_manager",
      "attributes": {
        "additional_sync_paths": [],
        "sync_disabled": false,
        "maximum_num_sync_threads": 10,
        "sync_interval_mins": 0.01,
        "capture_dir": "",
        "tags": [
          "nick"
        ]
      }
    }
  ],
  "modules": [
    {
      "type": "local",
      "name": "boomsensor",
      "executable_path": "/home/user/boomsensor"
    }
  ]
}
```
