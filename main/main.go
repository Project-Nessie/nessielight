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

	nessielight.InitDBwithFile("test.db")
	nessielight.InitV2rayService(inboundTag, vmessPort, vmessAddress, wsPath, v2rayApi)

	server.Register("/hello", "Hello!", nil, nil, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		if from == nil {
			server.Sendf(chatid, "invalid interaction")
			return
		}
		user, _ := GetUserByTid(from.ID)
		server.Sendf(chatid, "Hello!\nYour ID: <code>%d</code>\nAdministration: <b>%v</b>\nRegistered: <b>%v</b>",
			from.ID, isAdmin(from.ID), user != nil)
	})

	registerAdminService(&server)
	registerProxyService(&server)
	registerLoginService(&server)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	logger = log.New(os.Stderr, "[main] ", log.LstdFlags|log.Lmsgprefix)
}
