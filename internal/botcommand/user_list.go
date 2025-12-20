package botcommand

import (
	"fmt"
	"game-controller-bot/internal/botcontext"
	"maps"
	"slices"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
	"gopkg.in/telebot.v4"
)

type UserListCommand struct {
	context *botcontext.BotContext
}

func NewUserListCommand(context *botcontext.BotContext) BotCommand {
	return &UserListCommand{
		context: context,
	}
}

func (c *UserListCommand) GetName() string {
	return "/user_list"
}

func (c *UserListCommand) GetDescription() string {
	return "Список пользователей компьютера"
}

func (c *UserListCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		processes, err := process.Processes()
		if err != nil {
			return err
		}

		usernamesMap := make(map[string]any)
		for _, proc := range processes {
			username, err := proc.Username()
			if err != nil {
				return err
			}

			usernamesMap[username] = nil
		}

		usernames := slices.SortedFunc(maps.Keys(usernamesMap), strings.Compare)

		message := strings.Builder{}
		message.WriteString("Список пользователей:\n")
		for _, username := range usernames {
			message.WriteString(fmt.Sprintf("%v\n", username))
		}
		return context.Send(message.String())
	}
}
