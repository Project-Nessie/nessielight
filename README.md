# Nessie Light

:warning: This project is still under rapid development, thus everything may change.

<a href="https://github.com/Project-Nessie/nessielight/blob/master/LICENSE"><img alt="GitHub license" src="https://img.shields.io/github/license/Project-Nessie/nessielight"></a> <a href="https://app.fossa.com/projects/git%2Bgithub.com%2FProject-Nessie%2Fnessielight?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2FProject-Nessie%2Fnessielight.svg?type=shield"/></a> <a href="https://goreportcard.com/report/github.com/Project-Nessie/nessielight"><img alt="GitHub license" src="https://goreportcard.com/badge/github.com/Project-Nessie/nessielight"></a>

## Overview

_Nessie Light_ is a proxy manager based on telegram bot written in go, supporting multiuser, permission control, traffic limitation and data statistics. It works with v2ray and (in the future) other network tools. Taking advantage of telegram's interaction, it preserves security, simplicity and functionality.

_Nessie Light_ is planning to support web interaction in the future, which typically acts as an entry.

## Getting Start

Before install _Nessie Light_, you first install v2ray and enable it as a system service. i. e. follow the official installation guide.

Since _Nessie Light_ use api to communicate with v2ray, you should enable api and statistics in v2ray (add following configuration):

```json
{
  "stats": {},
  "api": {
    "tag": "api",
    "services": ["StatsService", "LoggerService", "HandlerService"]
  },
  "inbounds": [
    {
      "tag": "api",
      "listen": "127.0.0.1",
      "port": 10085,
      "protocol": "dokodemo-door",
      "settings": {
        "network": "tcp,udp",
        "address": "127.0.0.1"
      }
    }
  ]
}
```

_Nessie Light_ place all user in a single taged inbound, thus add an inbound with specified tag:

```json
{
  "inbounds": [
    {
      "tag": "multiuser",
      "listen": "0.0.0.0",
      "port": 38888,
      "protocol": "vmess",
      "settings": {
        "clients": []
      },
      "streamSettings": {
        "network": "ws",
        "wsSettings": {
          "path": "/apath"
        }
      },
      "sniffing": {
        "enabled": true,
        "destOverride": ["http", "tls"]
      }
    }
  ]
}
```

Now restart v2ray service.

Currently, you're supposed to build _Nessie Light_ from source to get the executable:

```bash
go build -trimpath -v -ldflags='-s -w -extldflags "-static"' -o dest/nessielight ./main
```

Usage:

```bash
$ ./nessielight --help
Usage of ./nessielight:
  -admin value
    	init admin using tg user id
  -listen string
    	listen address (default "127.0.0.1:3456")
  -token string
    	tg bot token
  -v2rayapi string
    	v2ray api listening address
  -vmessaddr string
    	vmess address
  -vmessport int
    	vmess port (default 443)
  -vmesstag string
    	vmess inbound tag
  -webhook string
    	tg bot webhook url
  -wspath string
    	websocket path
```

Lastly, create a telegram bot with webhook setting to corresponding url, and you're ready to start:

```bash
$ ./nessielight -token xxxxxx\
    -webhook https://xxxxxx:12345\
    -listen 0.0.0.0:12345\
    -vmessaddr example.com\
    -vmessport 38888\
    -vmesstag multiuser\
    -wspath /apath\
    -v2rayapi 127.0.0.1:10085\
```
