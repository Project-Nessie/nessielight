// package main implements command line interface
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/yanzay/tbot/v2"
)

var logger *log.Logger

func main() {
	flag.Parse()

	logger.Printf("Hello World!")
	server := tgolf.NewServer(botToken, webhookUrl, listenAddr)

	server.Register("/hello", "Hello!", nil, nil, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		server.Sendf(chatid, "Hello! (current: %v)", time.Now())
	})

	registerAdminService(&server)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	logger = log.New(os.Stderr, "[main] ", log.LstdFlags|log.Lmsgprefix)
}
