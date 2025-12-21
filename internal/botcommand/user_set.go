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
			return context.Send("Ожидается один параметр - имя пользователя операционной системы")
		}

		data, err := c.context.Read()
		if err != nil {
			return err
		}

		user := context.Args()[0]
		data.SetSystemUser(user)
		err = c.context.Write(data)
		if err != nil {
			return err
		}

		return context.Send(fmt.Sprintf("Пользователь '%v' установлен", user))
	}
}
