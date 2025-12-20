package botcommand

import (
	"fmt"
	"game-controller-bot/internal/botcontext"

	"gopkg.in/telebot.v4"
)

type UserSetCommand struct {
	context *botcontext.BotContext
}

func NewUserSetCommand(context *botcontext.BotContext) BotCommand {
	return &UserSetCommand{
		context: context,
	}
}

func (c *UserSetCommand) GetName() string {
	return "/user_set"
}

func (c *UserSetCommand) GetDescription() string {
	return "Задать имя пользователя компьютера, чьи процессы наблюдать"
}

func (c *UserSetCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if len(context.Args()) != 1 {
			return context.Send("Ожидается один параметр - имя пользователя")
		}

		user := context.Args()[0]
		c.context.SystemUser = user
		return context.Send(fmt.Sprintf("Пользователь '%v' установлен", user))
	}
}
