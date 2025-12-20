package botcommand

import (
	"fmt"
	"game-controller-bot/internal/botcontext"
	"strings"

	"gopkg.in/telebot.v4"
)

type FilterListCommand struct {
	context *botcontext.BotContext
}

func NewFilterListCommand(context *botcontext.BotContext) BotCommand {
	return &FilterListCommand{
		context: context,
	}
}

func (c *FilterListCommand) GetName() string {
	return "/filter_list"
}

func (c *FilterListCommand) GetDescription() string {
	return "Показывает текущие настройки фильтра"
}

func (c *FilterListCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if len(c.context.NamesFilter) == 0 {
			return context.Send("Фильтр пуст")
		}

		message := strings.Builder{}
		message.WriteString("Текущий фильтр:\n")
		for _, filter := range c.context.NamesFilter {
			message.WriteString(fmt.Sprintf("%v\n", filter))
		}
		return context.Send(message.String())

	}
}
