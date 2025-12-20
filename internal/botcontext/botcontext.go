package botcontext

import (
	"gopkg.in/telebot.v4"
)

type BotContext struct {
	Admin       *telebot.User
	SystemUser  string
	NamesFilter []string
}

func NewBotContext() *BotContext {
	return &BotContext{
		Admin:       nil,
		SystemUser:  "",
		NamesFilter: nil,
	}
}
