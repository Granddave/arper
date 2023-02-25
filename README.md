<div align="center">

![arper logo](res/arper.png)

<i>Simple Arp listener written in Go.</i>

</div>

[![Go CI](https://github.com/Granddave/arper/actions/workflows/go.yml/badge.svg)](https://github.com/Granddave/arper/actions/workflows/go.yml)

## About

`arper` is a tool that listens for Arp response packets and keeps track on the
devices on the network.

### Features

- Notification support (Discord Webhook API)
- Hostname lookup
- Hardware vendor lookup
- Persistant storage


## Usage

Root privileges are required since `arper` listen on raw packets on the specified
network interface.

```sh
sudo arper [flags]
```

### Configuration

Configurations can be set in multiple ways following the precedence order below:

1. Command line flag
2. Environment variable
3. Config file
4. Default value

```
      --database string          path to the database file (default "/var/lib/arper/database.json")
      --discord-webhook string   Discord webhook URL
  -h, --help                     help for arper
      --iface string             network interface to use (default "eth0")
```

## Build Requirements

Go version >= 1.18


## Build Instructions

```bash
go mod tidy
go build ./cmd/arper
```


## Roadmap

- Configuration file support
- Hardware vendor lookup caching to minimize the API requests


## Similar applications

- [arpwatch](https://linux.die.net/man/8/arpwatch)
