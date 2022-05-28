package main

import (
	"flag"
	"fmt"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprint([]string(*i))
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// telegram bot token https://core.telegram.org/bots#6-botfather
var botToken string

// webhook url for telegram bot https://core.telegram.org/bots/api#setwebhook
var webhookUrl string

// telegram bot server listening address
var listenAddr string

// v2ray api server listening address https://guide.v2fly.org/en_US/advanced/traffic.html#configuration-example
var v2rayApi string

// telegram user ID of admins
var admins arrayFlags

// tag of preconfigured inbound for user proxy
var inboundTag string

// vmess listening address
var vmessAddress string

// websocket path
var wsPath string

// vmess listening port
var vmessPort int

// vmess port connected by client (usually 443)
var vmessClientPort int

func init() {
	flag.StringVar(&botToken, "token", "", "tg bot token")
	flag.StringVar(&webhookUrl, "webhook", "", "tg bot webhook url")
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:3456", "listen address")
	flag.Var(&admins, "admin", "init admin using tg user id")
	flag.StringVar(&v2rayApi, "v2rayapi", "", "v2ray api listening address")
	flag.StringVar(&inboundTag, "vmesstag", "", "vmess inbound tag")
	flag.IntVar(&vmessClientPort, "vmessclientport", 443, "vmess client port")
	flag.IntVar(&vmessPort, "vmessport", 12345, "vmess listening port")
	flag.StringVar(&vmessAddress, "vmessaddr", "", "vmess address")
	flag.StringVar(&wsPath, "wspath", "", "websocket path")
}
