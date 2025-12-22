package botcommand

import (
	"fmt"
	"game-controller-bot/internal/botcontext"
	"game-controller-bot/internal/proc"
	"strings"

	"gopkg.in/telebot.v4"
)

type ProcessListCommand struct {
	context *botcontext.BotContext
}

func NewProcessListCommand(context *botcontext.BotContext) BotCommand {
	return &ProcessListCommand{
		context: context,
	}
}

func (c *ProcessListCommand) GetName() string {
	return "/process_list"
}

func (c *ProcessListCommand) GetDescription() string {
	return "Показывает запущенные процессы, которые подходят под фильтр"
}

func (c *ProcessListCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		data, err := c.context.Read()
		if err != nil {
			return err
		}

		grep := ""
		if len(context.Args()) == 1 {
			grep = strings.ToLower(context.Args()[0])
		}

		allProcesses, err := proc.ReadAllProcesses(nil, grep, data.SystemUser)
		if err != nil {
			return err
		}

		message := strings.Builder{}
		message.WriteString("Список всех процессов, подходящих под фильтр:\n")
		for _, p := range allProcesses {
			line := fmt.Sprintf("%v (> %v) %v\n", p.Pid, p.ParentPid, p.Name)
			if message.Len()+len(line) > 4096 {
				err = context.Send(message.String())
				if err != nil {
					return err
				}
				message.Reset()
			}
			message.WriteString(line)
		}
		return context.Send(message.String())
	}
}
