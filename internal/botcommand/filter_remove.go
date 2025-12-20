package botcommand

import (
	"game-controller-bot/internal/botcontext"
	"slices"

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
		for _, arg := range context.Args() {
			index := slices.Index(c.context.NamesFilter, arg)
			if index >= 0 {
				c.context.NamesFilter = slices.Delete(c.context.NamesFilter, index, index+1)
			}
		}

		return NewUserListCommand(c.context).GetHandleFunc()(context)
	}
}
