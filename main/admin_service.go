package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Project-Nessie/nessielight"
	"github.com/Project-Nessie/nessielight/tgolf"
	"github.com/Project-Nessie/nessielight/utils"
	"github.com/yanzay/tbot/v2"
)

// provide admin using cli
func isAdmin(userid int) bool {
	id := fmt.Sprint(userid)
	for _, v := range admins {
		if id == v {
			return true
		}
	}
	return false
}

var adminHelpText = `
`

var userManHelp = `
Add User: generate new token for registering
`

func registerAdminService(server *tgolf.Server) {
	adminBtns := [][]tbot.InlineKeyboardButton{
		{{Text: "User Management", CallbackData: "a/user"}},
		{{Text: "Service Control", CallbackData: "a/service"}},
		{{Text: "Statistics", CallbackData: "a/statistics"}},
	}
	userManBtns := [][]tbot.InlineKeyboardButton{
		{{Text: "Add User", CallbackData: "a/user/add"}, {Text: "Delete User", CallbackData: "a/user/delete"}},
		{{Text: "Set User", CallbackData: "a/user/set"}, {Text: "Go Back", CallbackData: "a/back"}},
	}
	serviceBtns := [][]tbot.InlineKeyboardButton{
		{{Text: "Restart V2ray", CallbackData: "a/service/v2rayrestart"}},
		{{Text: "View V2ray Log", CallbackData: "a/service/v2raylog"}},
		{{Text: "Go Back", CallbackData: "a/back"}},
	}
	statisBtns := [][]tbot.InlineKeyboardButton{
		{{Text: "Get Top Traffic", CallbackData: "a/statistics/toptraffic"}},
		{{Text: "Get Top Traffic Today", CallbackData: "a/statistics/toptraffictoday"}},
		{{Text: "Reset Traffic", CallbackData: "a/statistics/resettraffic"}},
		{{Text: "Go Back", CallbackData: "a/back"}},
	}

	server.Register("/admin", "Admin Control", combineInit(withPrivate, withAdmin), nil,
		func(argv []tgolf.Argument, from *tbot.User, chatid string) {
			server.SendfWithBtn(chatid, adminBtns, "Your User ID: %d\n%s", from.ID, adminHelpText)
		})

	server.RegisterInlineButton("a/back", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsgWithBtn(cq, adminBtns, "Your User ID: %d\n%s", cq.From.ID, adminHelpText)
		return nil
	})
	server.RegisterInlineButton("a/user", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsgWithBtn(cq, userManBtns, "User Management\n%s", userManHelp)
		return nil
	})
	// 生成一个 token，用于注册用户
	server.RegisterInlineButton("a/user/add", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackBtn(cq, [][]tbot.InlineKeyboardButton{})
		token := nessielight.AuthServiceInstance.GenToken()
		server.Sendf(cq.Message.Chat.ID, "token: <code>%s</code>", token)
		return nil
	})

	server.Register(">>>user/delete", "", nil, []tgolf.Parameter{
		tgolf.NewParam("id", "user id", func(value string) bool {
			id, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return false
			}
			_, err = GetUserByTid(int(id))
			return err == nil
		}),
	}, func(argv []tgolf.Argument, from *tbot.User, chatid string) {
		id, err := strconv.ParseInt(argv[0].Value, 10, 64)
		if err != nil {
			logger.Print(err)
			return
		}
		user, err := GetUserByTid(int(id))
		if err != nil {
			logger.Print(err)
			return
		}
		if err := nessielight.UserManagerInstance.DeleteUser(user); err != nil {
			logger.Print(err)
			return
		}
		server.Sendf(chatid, "done.")
	})

	server.RegisterInlineButton("a/user/delete", func(cq *tbot.CallbackQuery) error {
		users, err := nessielight.UserManagerInstance.All()
		if err != nil {
			return err
		}
		msg := "Users:\n"
		for _, v := range users {
			msg += v.Name() + ": <code>" + v.ID() + "</code>" + "\n"
		}
		server.EditCallbackBtn(cq, [][]tbot.InlineKeyboardButton{})
		server.Sendf(cq.Message.Chat.ID, msg)
		if err := server.StartCommand(">>>user/delete", cq.From, cq.Message.Chat); err != nil {
			return err
		}
		return nil
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/user/set", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsg(cq, "<i>set user not implemented</i>")
		return nil
	})

	server.RegisterInlineButton("a/service", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsgWithBtn(cq, serviceBtns, "Service Control")
		return nil
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/service/v2rayrestart", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsg(cq, "<i>v2ray start not implemented</i>")
		return nil
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/service/v2raylog", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsg(cq, "<i>v2ray log not implemented</i>")
		return nil
	})

	server.RegisterInlineButton("a/statistics", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsgWithBtn(cq, statisBtns, "Service Control")
		return nil
	})
	server.RegisterInlineButton("a/statistics/toptraffic", func(cq *tbot.CallbackQuery) error {
		stats, err := nessielight.GetV2rayTraffic()
		if err != nil {
			return err
		}
		users, inbounds := stats.Users(), stats.Inbounds()
		sort.Sort(users)
		sort.Sort(inbounds)

		parseValue := func(num int64) string {
			if num < 1000 {
				return fmt.Sprint(num, "B")
			} else if num < 1e6 {
				return fmt.Sprint(int(num/1e3), "KB")
			} else if num < 1e9 {
				return fmt.Sprint(int(num/1e6), "MB")
			} else {
				return fmt.Sprint(int(num/1e9), "GB")
			}
		}

		reducer := func(msg string, v nessielight.UserTraffic) string {
			name := v.Name
			if v.Type == nessielight.Unregistered {
				name = "<i>" + v.Name + "</i>"
			}
			return fmt.Sprintf("%s%s down <b>%s</b> up <b>%s</b>\n", msg, name,
				parseValue(v.Downlink), parseValue(v.Uplink))
		}

		msg := "<b><u>User Traffics sorted by downlink</u></b>\n"
		msg = utils.Reduce(users, reducer, msg)
		msg = msg + "\n<b><u>Inbound Traffics sorted by downlink</u></b>\n"
		msg = utils.Reduce(inbounds, reducer, msg)

		server.EditCallbackMsg(cq, msg)
		return nil
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/statistics/toptraffictoday", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsg(cq, "<i>top traffic today not implemented</i>")
		return nil
	})
	// !!!UNIMPLEMENTED
	server.RegisterInlineButton("a/statistics/resettraffic", func(cq *tbot.CallbackQuery) error {
		server.EditCallbackMsg(cq, "<i>reset traffic not implemented</i>")
		return nil
	})
}
