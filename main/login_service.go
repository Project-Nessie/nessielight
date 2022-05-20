package main

import (
	"fmt"

	"github.com/Project-Nessie/nessielight"
	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/yanzay/tbot/v2"
)

func registerLoginService(server *tgolf.Server) {
	server.Register("/register", "Register yourself", combineInit(withPrivate, func(from *tbot.User, chat tbot.Chat) bool {
		id := fmt.Sprint(from.ID)
		user, err := nessielight.UserManagerInstance.FindUser(id)
		if err != nil {
			server.Sendf(chat.ID, err.Error())
			return false
		}
		if user != nil {
			server.Sendf(chat.ID, "You've already registered")
			return false
		}
		return true
	}), []tgolf.Parameter{
		tgolf.NewParam("token", "token", nil),
	}, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		token := argv[0].Value
		user, err := nessielight.AuthServiceInstance.Register(token, fmt.Sprint(from.ID))
		if err != nil {
			server.Sendf(chatid, "register failed: %s", err.Error())
			return
		}
		if err := user.SetName(from.Username); err != nil {
			server.Sendf(chatid, "register failed: %s", err.Error())
			return
		}
		if err := nessielight.UserManagerInstance.SetUser(user); err != nil {
			server.Sendf(chatid, "register failed: %s", err.Error())
			return
		}
		server.Sendf(chatid, "Register succeed.")
	})
}
