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

Configurations can be set in multiple ways following the precedence order below.

1. Command line flag
2. Config file
3. Default value

arper doesn't create any configuration file, but looks for `./arper.yaml` first and
`/etc/arper/arper.yaml` second.

| Configuration | Config file key | Default | CLI Flag |
|-|-|-|-|-|
| Listening interface | `Iface`             | `eth0`                         | `--iface [inteface]`          |
| Database filepath   | `DatabaseFilepath`  | `/var/lib/arper/database.json` | `--database-filepath [path]`  |
| Discord webhook URL | `DiscordWebhookURL` | *(empty)*                      | `--discord-webhook-url [url]` |


## Build Requirements

Go version >= 1.18


## Build Instructions

```bash
go mod tidy
go build ./cmd/arper
```


## Similar applications

- [arpwatch](https://linux.die.net/man/8/arpwatch)
