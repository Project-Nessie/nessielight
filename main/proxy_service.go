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

	server.RegisterInlineButton("p/back", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsgWithBtn(cq, proxyBtns, "<b>Proxy Control</b>\nYour User ID: %d", cq.From.ID)
		return nil
	})
	server.RegisterInlineButton("p/get", func(cq *tbot.CallbackQuery) error {
		id := fmt.Sprint(cq.From.ID)
		user, err := nessielight.UserManagerInstance.FindUser(id)
		if err != nil {
			return err
		}
		nessielight.V2rayServiceInstance.RemoveUser(user.ID())
		uuid, err := nessielight.V2rayServiceInstance.AddUser(user.ID())
		if err != nil {
			return err
		}
		server.Sendf(cq.Message.Chat.ID, "uuid: %s", uuid)
		logger.Printf("uuid: %s", uuid)
		return nil
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("p/upd", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsg(cq, "<i>update configs not implemented</i>")
		return nil
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("p/stat", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsg(cq, "<i>traffic statistics not implemented</i>")
		return nil
	})
}
