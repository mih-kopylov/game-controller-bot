package botcommand

import (
	"gopkg.in/telebot.v4"
)

type BotCommand interface {
	GetName() string
	GetDescription() string
	GetHandleFunc() telebot.HandlerFunc
}
