package botcommand

import (
	"game-controller-bot/internal/botcontext"

	"gopkg.in/telebot.v4"
)

type FilterAddCommand struct {
	context *botcontext.BotContext
}

func NewFilterAddCommand(context *botcontext.BotContext) BotCommand {
	return &FilterAddCommand{
		context: context,
	}
}

func (c *FilterAddCommand) GetName() string {
	return "/filter_add"
}

func (c *FilterAddCommand) GetDescription() string {
	return "Добавляет одно или несколько значений в фильтр"
}

func (c *FilterAddCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		data, err := c.context.Read()
		if err != nil {
			return err
		}

		data.AddNames(context.Args())
		err = c.context.Write(data)
		if err != nil {
			return err
		}

		return NewFilterListCommand(c.context).GetHandleFunc()(context)
	}
}
