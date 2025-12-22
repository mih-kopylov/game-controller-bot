package botcommand

import (
	"game-controller-bot/internal/botcontext"

	"gopkg.in/telebot.v4"
)

type FilterRemoveCommand struct {
	context *botcontext.BotContext
}

func NewFilterRemoveCommand(context *botcontext.BotContext) BotCommand {
	return &FilterRemoveCommand{
		context: context,
	}
}

func (c *FilterRemoveCommand) GetName() string {
	return "/filter_remove"
}

func (c *FilterRemoveCommand) GetDescription() string {
	return "Удаляет одно или несколько значений из фильтра"
}

func (c *FilterRemoveCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		data, err := c.context.Read()
		if err != nil {
			return err
		}

		data.RemoveNames(context.Args())
		err = c.context.Write(data)
		if err != nil {
			return err
		}

		return NewFilterListCommand(c.context).GetHandleFunc()(context)
	}
}
