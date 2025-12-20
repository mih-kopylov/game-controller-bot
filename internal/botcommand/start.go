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
		admin := context.Sender()
		c.context.Admin = admin
		err := context.Send(fmt.Sprintf("Hello, %v (%v)", admin.Username, admin.ID))
		if err != nil {
			return err
		}

		return nil
	}
}
