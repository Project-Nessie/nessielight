package main

import (
	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/yanzay/tbot/v2"
)

func registerLoginService(server *tgolf.Server) {
	// !!!UNIMPLEMENTED
	server.Register("/register", "Register yourself", func(from *tbot.User, chat tbot.Chat) bool {
		return from != nil && chat.Type == "private"
	}, nil, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		server.Sendf(chatid, "<i>register service not implemented</i>")
	})
}
