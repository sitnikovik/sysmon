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

## API

There is a gRPC API to get the app results stored in `tmp/`

### GetStats

Returns the system monitoring statistics resulted with the app.

#### Response example

```json
{
    "cpu": {
        "user": 8.12,
        "system": 11.5,
        "idle": 80.82
    },
    "disk": {
        "reads": 0,
        "writes": 0,
        "readWriteKB": 0,
        "totalMb": "0",
        "usedMb": "0",
        "usedPercent": 0,
        "usedInodes": "0",
        "usedInodesPercent": 0
    },
    "memory": {
        "totalMb": "21181",
        "availableMb": "9390",
        "freeMb": "202",
        "activeMb": "9233",
        "inactiveMb": "9187",
        "WiredMb": "2551"
    },
    "loadAverage": {
        "oneMin": 2.32,
        "fiveMin": 2.47,
        "fifteenMin": 2.68
    }
}
```
