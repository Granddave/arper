<div align="center">

![arper logo](res/arper.png)

<i>Simple Arp listener written in Go.</i>

</div>

[![Go CI](https://github.com/Granddave/arper/actions/workflows/go.yml/badge.svg)](https://github.com/Granddave/arper/actions/workflows/go.yml)

## About

`arper` is a tool that listens for Arp response packets and keeps track on the
devices on the network. New hosts are discovered and stored in a local database.


## Usage

Root privileges are required since `arper` listen on raw packets on the specified
network interface.

```bash
sudo arper [-db PATH] [-iface IFACE] [-discord-webhook URL]
```

### Command line arguments

```bash
  -db string
        filepath to database (default "/var/lib/arper/hosts.json")
  -discord-webhook string
        Discord Webhook URL for notifications
  -iface string
        network interface to use (default "eth0")
```


## Build Requirements

Go version >= 1.18


## Build Instructions

```bash
go build
```


## Roadmap

- Configuration file support
- Hardware vendor lookup

## Similar applications

- [arpwatch](https://linux.die.net/man/8/arpwatch)
