package main

import (
	"fmt"

	"github.com/Project-Nessie/nessielight/service"
	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/yanzay/tbot/v2"
)

func registerLoginService(server *tgolf.Server) {
	server.Register("/register", "Register yourself", func(from *tbot.User, chat tbot.Chat) bool {
		if from == nil || chat.Type != "private" {
			server.Sendf(chat.ID, "invalid environment")
			return false
		}
		id := fmt.Sprint(from.ID)
		user, err := service.UserManagerInstance.FindUser(id)
		if err != nil {
			server.Sendf(chat.ID, err.Error())
			return false
		}
		if user != nil {
			server.Sendf(chat.ID, "You've already registered")
			return false
		}
		return true
	}, []tgolf.Parameter{
		tgolf.NewParam("token", "token", nil),
	}, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		token := argv[0].Value
		if _, err := service.AuthServiceInstance.Register(token, fmt.Sprint(from.ID)); err != nil {
			server.Sendf(chatid, "register failed: %s", err.Error())
			return
		}
		server.Sendf(chatid, "Register succeed.")
	})
}
