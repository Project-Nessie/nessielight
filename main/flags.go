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

var botToken, webhookUrl, listenAddr, v2rayApi string
var admins arrayFlags

// flags
var inboundTag, vmessAddress, wsPath string
var vmessPort int

func init() {
	flag.StringVar(&botToken, "token", "", "tg bot token")
	flag.StringVar(&webhookUrl, "webhook", "", "tg bot webhook url")
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:3456", "listen address")
	flag.Var(&admins, "admin", "init admin using tg user id")
	flag.StringVar(&v2rayApi, "v2rayapi", "", "v2ray api listening address")
	flag.StringVar(&inboundTag, "vmesstag", "", "vmess inbound tag")
	flag.IntVar(&vmessPort, "vmessport", 443, "vmess port")
	flag.StringVar(&vmessAddress, "vmessaddr", "", "vmess address")
	flag.StringVar(&wsPath, "wspath", "", "websocket path")
}
