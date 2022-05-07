package main

import (
	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/yanzay/tbot/v2"
)

// !!!UNIMPLEMENTED
func isAdmin(userid int) bool {
	return true
}

func registerAdminService(server *tgolf.Server) {
	adminBtns := [][]tbot.InlineKeyboardButton{
		{{
			Text:         "User Management",
			CallbackData: "a/user",
		}}, {{
			Text:         "Service Control",
			CallbackData: "a/service",
		}}, {{
			Text:         "Statistics",
			CallbackData: "a/statistics",
		}},
	}
	userManBtns := [][]tbot.InlineKeyboardButton{
		{{
			Text:         "Add User",
			CallbackData: "a/user/add",
		}}, {{
			Text:         "Delete User",
			CallbackData: "a/user/delete",
		}}, {{
			Text:         "Set User",
			CallbackData: "a/user/set",
		}}, {{
			Text:         "Go Back",
			CallbackData: "a/back",
		}},
	}
	serviceBtns := [][]tbot.InlineKeyboardButton{
		{{
			Text:         "Restart V2ray",
			CallbackData: "a/service/v2rayrestart",
		}}, {{
			Text:         "View V2ray Log",
			CallbackData: "a/service/v2raylog",
		}}, {{
			Text:         "Go Back",
			CallbackData: "a/back",
		}},
	}
	statisBtns := [][]tbot.InlineKeyboardButton{
		{{
			Text:         "Get Top Traffic",
			CallbackData: "a/statistics/toptraffic",
		}}, {{
			Text:         "Get Top Traffic Today",
			CallbackData: "a/statistics/toptraffictoday",
		}}, {{
			Text:         "Reset Traffic",
			CallbackData: "a/statistics/resettraffic",
		}}, {{
			Text:         "Go Back",
			CallbackData: "a/back",
		}},
	}

	server.Register("/admin", "Admin Control", func(from *tbot.User, chat tbot.Chat) bool {
		return from != nil && isAdmin(from.ID)
	}, nil, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		server.SendfWithBtn(chatid, adminBtns, "Your User ID: %d", from.ID)
	})

	server.RegisterInlineButton("a/back", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsgWithBtn(cq, adminBtns, "Your User ID: %d", cq.From.ID)
	})
	server.RegisterInlineButton("a/user", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsgWithBtn(cq, userManBtns, "User Management")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/user/add", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>add user not implemented</i>")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/user/delete", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>delete user not implemented</i>")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/user/set", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>set user not implemented</i>")
	})

	server.RegisterInlineButton("a/service", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsgWithBtn(cq, serviceBtns, "Service Control")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/service/v2rayrestart", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>v2ray start not implemented</i>")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/service/v2raylog", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>v2ray log not implemented</i>")
	})

	server.RegisterInlineButton("a/statistics", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsgWithBtn(cq, statisBtns, "Service Control")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/statistics/toptraffic", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>top traffic not implemented</i>")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/statistics/toptraffictoday", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>top traffic today not implemented</i>")
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/statistics/resettraffic", func(cq *tbot.CallbackQuery) {
		server.EditCallbackMsg(cq, "<i>reset traffic not implemented</i>")
	})
}
