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
		for _, arg := range context.Args() {
			c.context.NamesFilter = append(c.context.NamesFilter, arg)
		}

		return NewUserListCommand(c.context).GetHandleFunc()(context)
	}
}
