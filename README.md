# sysmon

[![Build Status](https://github.com/sitnikovik/sysmon/actions/workflows/go.yml/badge.svg)](https://github.com/sitnikovik/sysmon/actions)
[![Coverage Status](https://coveralls.io/repos/github/sitnikovik/sysmon/badge.svg?branch=master)](https://coveralls.io/github/sitnikovik/sysmon?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sitnikovik/sysmon)](https://goreportcard.com/report/github.com/sitnikovik/sysmon)
[![License](https://img.shields.io/github/license/sitnikovik/sysmon)](https://github.com/sitnikovik/sysmon/blob/master/LICENSE)

Daemon program to collect information about the system is running on and sends it to its clients in GRPC.

## Supported platforms

- OS X (Tested on 14.6.1)
- Ubuntu 18.04

## Getting started

- Clone the repo
- Install the dependencies

```sh
make install-deps 
```

- Build the app

```sh
# Creates the binary bin/sysmon
make build
```

## Usage

- `-n` - interval of time to output the metrics
- `-m` - margin of time between statistics output
- `-grpc-port` - gRPC port to run the gRPC-server to get metrics by API
- `--config` - path to the configuration yaml-file that stores all app settings

> Config values replaces flag values

```sh
# Runs the app with yaml configuration file
bin/sysmon --confifg=path/to/sysmon.yml
```

```sh
# Runs the app with manul provided flags
bin/sysmon -n=5 -m=15 -grpc-port=50051
```

```sh
# Runs the app with only default configuration
bin/sysmon
```

### Example

[![Output example](output_example.png)]


## API

There is a gRPC API to get the app results stored in `tmp/`

### GetStats

Returns the system monitoring statistics resulted with the app.

#### Response example

```json
{
    "cpu": {
        "user": 7.48,
        "system": 11.2,
        "idle": 81.49
    },
    "disk": {
        "reads": 34,
        "writes": 0,
        "readWriteKb": 349.85999999999996,
        "totalMb": "1018880",
        "usedMb": "745472",
        "usedPercent": 76,
        "usedInodes": "1422283552",
        "usedInodesPercent": 76
    },
    "memory": {
        "totalMb": "19626",
        "availableMb": "8427",
        "freeMb": "154",
        "activeMb": "8330",
        "inactiveMb": "8272",
        "WiredMb": "2817"
    },
    "loadAverage": {
        "oneMin": 2.46,
        "fiveMin": 2.97,
        "fifteenMin": 3.13
    }
}
```
