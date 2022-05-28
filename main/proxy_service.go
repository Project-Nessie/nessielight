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
		user, err := GetUserByTid(cq.From.ID)
		if err != nil {
			return err
		}
		nessielight.ApplyUserProxy(user)
		server.Sendf(cq.Message.Chat.ID, nessielight.GetUserProxyMessage(user))
		return nil
	})
	server.RegisterInlineButton("p/upd", func(cq *tbot.CallbackQuery) error {
		if err := nessielight.V2rayUpdateUserTraffic(); err != nil {
			return err
		}
		user, err := GetUserByTid(cq.From.ID)
		if err != nil {
			return err
		}
		for _, p := range user.Proxy() {
			p.Deactivate()
		}
		proxy := nessielight.V2rayServiceInstance.NewProxy()
		if err := user.SetProxy([]nessielight.Proxy{proxy}); err != nil {
			return err
		}
		if err := nessielight.ApplyUserProxy(user); err != nil {
			return err
		}
		if err := nessielight.UserManagerInstance.SetUser(user); err != nil {
			return err
		}
		server.EditCallbackMsgWithBtn(cq, [][]tbot.InlineKeyboardButton{}, "Proxy has updated.")
		server.Sendf(cq.Message.Chat.ID, nessielight.GetUserProxyMessage(user))
		return nil
	})

	server.RegisterInlineButton("p/stat", func(cq *tbot.CallbackQuery) error {
		user, err := GetUserByTid(cq.From.ID)
		if err != nil {
			return err
		}
		if err := nessielight.V2rayUpdateUserTraffic(); err != nil {
			return err
		}
		traffic := user.Traffic()
		server.EditCallbackMsgWithBtn(cq, [][]tbot.InlineKeyboardButton{},
			fmt.Sprintf("down <b>%v</b> up <b>%v</b>\n", traffic.Downlink, traffic.Uplink))
		return nil
	})
}
