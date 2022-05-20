package main

import (
	"fmt"

	"github.com/Project-Nessie/nessielight"
	"github.com/yanzay/tbot/v2"
)

func withPrivate(from *tbot.User, chat tbot.Chat) bool {
	return from != nil && chat.Type == "private"
}

func withAdmin(from *tbot.User, chat tbot.Chat) bool {
	return from != nil && isAdmin(from.ID)
}

func withAuth(from *tbot.User, chat tbot.Chat) bool {
	id := fmt.Sprint(from.ID)
	user, err := nessielight.UserManagerInstance.FindUser(id)
	if err != nil {
		logger.Print("error: ", err)
		return false
	}
	if user == nil { // haven't registered
		return false
	}
	return true
}

func combineInit(inits ...func(from *tbot.User, chat tbot.Chat) bool) func(*tbot.User, tbot.Chat) bool {
	return func(u *tbot.User, c tbot.Chat) bool {
		for _, v := range inits {
			if !v(u, c) {
				return false
			}
		}
		return true
	}
}

func GetUserByTid(id int) (nessielight.User, error) {
	uid := fmt.Sprint(id)
	user, err := nessielight.UserManagerInstance.FindUser(uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
