package botcommand

import (
	"game-controller-bot/internal/botcontext"

	"gopkg.in/telebot.v4"
)

type UserClearCommand struct {
	context *botcontext.BotContext
}

func NewUserClearCommand(context *botcontext.BotContext) BotCommand {
	return &UserClearCommand{
		context: context,
	}
}

func (c *UserClearCommand) GetName() string {
	return "/user_clear"
}

func (c *UserClearCommand) GetDescription() string {
	return "Очистить имя пользователя компьютера, чьи процессы наблюдать. Бот будет наблюдать всех пользователей"
}

func (c *UserClearCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		data, err := c.context.Read()
		if err != nil {
			return err
		}

		data.ClearSystemUser()
		err = c.context.Write(data)
		if err != nil {
			return err
		}

		return context.Send("Пользователь очищен")
	}
}
