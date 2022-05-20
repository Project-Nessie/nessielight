package main

import (
	"fmt"

	"github.com/Project-Nessie/nessielight"
	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/yanzay/tbot/v2"
)

func registerProxyService(server *tgolf.Server) {
	proxyBtns := [][]tbot.InlineKeyboardButton{
		{{Text: "Get Configs", CallbackData: "p/get"}},
		{{Text: "Update Configs", CallbackData: "p/upd"}},
		{{Text: "Get Statistics", CallbackData: "p/stat"}},
	}
	server.Register("/proxy", "Proxy Control", combineInit(withPrivate, withAuth), nil,
		func(argv []tgolf.Argument, from *tbot.User, chatid string) {
			server.SendfWithBtn(chatid, proxyBtns, "<b>Proxy Control</b>\nYour User ID: %d", from.ID)
		})

	server.RegisterInlineButton("p/back", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsgWithBtn(cq, proxyBtns, "<b>Proxy Control</b>\nYour User ID: %d", cq.From.ID)
	})
	server.RegisterInlineButton("p/get", func(cq *tbot.CallbackQuery) {
		id := fmt.Sprint(cq.From.ID)
		user, err := nessielight.UserManagerInstance.FindUser(id)
		if err != nil {
			server.Sendf(cq.Message.Chat.ID, "[finduser]: %s", err.Error())
			logger.Print(err)
			return
		}
		nessielight.V2rayServiceInstance.RemoveUser(user.ID())
		uuid, err := nessielight.V2rayServiceInstance.AddUser(user.ID())
		if err != nil {
			server.Sendf(cq.Message.Chat.ID, "[adduser]: %s", err.Error())
			logger.Print(err)
			return
		}
		server.Sendf(cq.Message.Chat.ID, "uuid: %s", uuid)
		logger.Printf("uuid: %s", uuid)
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("p/upd", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>update configs not implemented</i>")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("p/stat", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>traffic statistics not implemented</i>")
	})
}
