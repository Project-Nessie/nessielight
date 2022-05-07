package main

import (
	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/yanzay/tbot/v2"
)

// !!!UNIMPLEMENTED
func isRegistered(userid int) bool {
	return true
}

func registerProxyService(server *tgolf.Server) {
	proxyBtns := [][]tbot.InlineKeyboardButton{
		{{
			Text:         "Get Configs",
			CallbackData: "p/get",
		}}, {{
			Text:         "Update Configs",
			CallbackData: "p/upd",
		}}, {{
			Text:         "Get Statistics",
			CallbackData: "p/stat",
		}},
	}
	server.Register("/proxy", "Proxy Control", func(from *tbot.User, chat tbot.Chat) bool {
		return from != nil && isRegistered(from.ID)
	}, nil, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		server.SendfWithBtn(chatid, proxyBtns, "<b>Proxy Control</b>\nYour User ID: %d", from.ID)
	})

	server.RegisterInlineButton("p/back", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsgWithBtn(cq, proxyBtns, "<b>Proxy Control</b>\nYour User ID: %d", cq.From.ID)
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("p/get", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>get configs not implemented</i>")
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
