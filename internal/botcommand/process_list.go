package botcommand

import (
	"fmt"
	"game-controller-bot/internal/botcontext"
	"slices"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
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

		allProcesses, err := c.readAllProcesses(data, grep)
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

func (c *ProcessListCommand) readAllProcesses(data *botcontext.Data, grep string) ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var allProcesses []ProcessInfo
	for _, p := range processes {
		processName, _ := p.Name()
		processExe, _ := p.Exe()
		processUsername, _ := p.Username()
		terminal, _ := p.Terminal()
		parentPid, _ := p.Ppid()
		if len(data.NamesToKill) > 0 && !slices.Contains(data.NamesToKill, strings.ToLower(processName)) {
			continue
		}
		if data.SystemUser != "" && processUsername != data.SystemUser {
			continue
		}
		if grep != "" && !strings.Contains(strings.ToLower(processName), grep) {
			continue
		}

		allProcesses = append(
			allProcesses, ProcessInfo{
				Pid:       p.Pid,
				Name:      processName,
				User:      processUsername,
				Exe:       processExe,
				ParentPid: parentPid,
				Terminal:  terminal,
			},
		)
	}

	slices.SortFunc(
		allProcesses, func(a, b ProcessInfo) int {
			if a.User != b.User {
				return strings.Compare(a.User, b.User)
			}
			if a.Name != b.Name {
				return strings.Compare(a.Name, b.Name)
			}
			return 0
		},
	)
	return allProcesses, nil
}

type ProcessInfo struct {
	Pid       int32
	Name      string
	User      string
	Exe       string
	ParentPid int32
	Terminal  string
}
