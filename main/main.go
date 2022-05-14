// package main implements command line interface
package main

import (
	"flag"
	"log"
	"os"

	"github.com/Project-Nessie/nessielight/service"
	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/yanzay/tbot/v2"
)

var logger *log.Logger

func main() {
	flag.Parse()

	logger.Printf("Hello World!")
	server := tgolf.NewServer(botToken, webhookUrl, listenAddr)

	server.Register("/hello", "Hello!", nil, nil, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		server.Sendf(chatid, "Hello!\nYour ID: <code>%d</code>\nAdministration: <b>%v</b>", from.ID, isAdmin(from.ID))
	})

	registerAdminService(&server)
	registerProxyService(&server)
	registerLoginService(&server)

	client := service.NewV2rayClient(inboundTag, vmessPort, vmessAddress, wsPath)
	service.V2rayServiceInstance = &client

	if err := service.V2rayServiceInstance.Start(v2rayApi); err != nil {
		log.Fatal(err)
	}

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	logger = log.New(os.Stderr, "[main] ", log.LstdFlags|log.Lmsgprefix)
}
