// package main implements command line interface
package main

import (
	"flag"
	"log"
	"os"

	"github.com/Project-Nessie/nessielight"
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

	nessielight.InitDBwithFile("test.db")
	nessielight.InitV2rayService(inboundTag, vmessPort, vmessAddress, wsPath, v2rayApi)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	logger = log.New(os.Stderr, "[main] ", log.LstdFlags|log.Lmsgprefix)
}
