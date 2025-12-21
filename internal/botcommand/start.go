package botcommand

import (
	"fmt"
	"game-controller-bot/internal/botcontext"

	"gopkg.in/telebot.v4"
)

type StartCommand struct {
	context *botcontext.BotContext
}

func NewStartCommand(context *botcontext.BotContext) BotCommand {
	return &StartCommand{
		context: context,
	}
}

func (c *StartCommand) GetName() string {
	return "/start"
}

func (c *StartCommand) GetDescription() string {
	return "Начинает работу с ботом"
}

func (c *StartCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		data, err := c.context.Read()
		if err != nil {
			return err
		}

		user := context.Sender()
		switch {
		case len(data.Admins) == 0:
			data.AddAdmin(user.ID)
			err = c.context.Write(data)
			if err != nil {
				return err
			}

			return context.Send(
				fmt.Sprintf(
					`Привет, %v! 
Ты - первый пользователь бота, теперь ты его администратор.
Управляй им мудро!`,
					user.Username,
				),
			)

		case data.IsAdmin(user.ID):
			return context.Send(
				fmt.Sprintf(
					`Привет, %v!
Ты - администратор этого бота.
Управляй им мудро!`, user.Username,
				),
			)

		default:
			return context.Send(
				fmt.Sprintf(
					`Привет, %v!
Ты - не администратор этого бота.
Для того, чтобы управлять им, попроси администратора выдать тебе права.
`, user.Username,
				),
			)

		}
	}
}
