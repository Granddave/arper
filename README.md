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

| Configuration | CLI Flag | Config file key | Default |
|---|---|---|---|
| Listening interface | `--iface [inteface]`          | `Iface`             | `eth0`                         |
| Database filepath   | `--database-filepath [path]`  | `DatabaseFilepath`  | `/var/lib/arper/database.json` |
| Discord webhook URL | `--discord-webhook-url [url]` | `DiscordWebhookURL` | *(empty)*                      |


## Build and run in Docker

```sh
docker build -t arper .
```

```sh
docker run \
  --rm -it \
  --net host \
  -v $(pwd)/database.json:/var/lib/arper/database.json \
  -v $(pwd)/arper.yaml:/arper.yaml \
  arper
```


## Local build

Go version >= 1.18 is required.

```bash
go mod tidy
go build ./cmd/arper
```


## Similar applications

- [arpwatch](https://linux.die.net/man/8/arpwatch)
