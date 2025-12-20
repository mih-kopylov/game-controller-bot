package botcommand

import (
	"fmt"
	"game-controller-bot/internal/botcontext"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
	"gopkg.in/telebot.v4"
)

type ProcessTerminateCommand struct {
	context *botcontext.BotContext
}

func NewProcessTerminateCommand(context *botcontext.BotContext) BotCommand {
	return &ProcessTerminateCommand{
		context: context,
	}
}

func (c *ProcessTerminateCommand) GetName() string {
	return "/process_terminate"
}

func (c *ProcessTerminateCommand) GetDescription() string {
	return "Останавливает процессы по переданным pid"
}

func (c *ProcessTerminateCommand) GetHandleFunc() telebot.HandlerFunc {
	return func(context telebot.Context) error {
		message := strings.Builder{}

		var pids []int32
		for _, pidString := range context.Args() {
			pid, err := strconv.Atoi(pidString)
			if err != nil {
				message.WriteString(fmt.Sprintf("Значение '%v' - не число\n", pidString))
				continue
			}
			pids = append(pids, int32(pid))
		}

		for _, pid := range pids {
			proc, err := process.NewProcess(pid)
			if err != nil {
				message.WriteString(fmt.Sprintf("Процесс с pid='%v' не найден\n", pid))
				continue
			}

			err = proc.Terminate()
			if err != nil {
				message.WriteString(fmt.Sprintf("Не удалось завершить процесс с pid='%v'\n", pid))
				continue
			}

			message.WriteString(fmt.Sprintf("Процесс с pid='%v' завершён\n", pid))
		}

		return context.Send(message.String())
	}
}
