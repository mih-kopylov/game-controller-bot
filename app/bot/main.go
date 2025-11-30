package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/process"
	"gopkg.in/telebot.v4"
)

func main() {
	pref := telebot.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	var admin *telebot.User
	var namesFilter []string

	bot.Handle(
		"/start", func(context telebot.Context) error {
			admin = context.Sender()
			return context.Send(fmt.Sprintf("Hello, %v (%v)", admin.Username, admin.ID))
		},
	)
	bot.Handle(
		"/hello", func(context telebot.Context) error {
			if admin == nil {
				return context.Send("I don't know you. Run /start to start")
			} else {
				return context.Send(fmt.Sprintf(":wave:, %v (%v)", admin.Username, admin.ID))
			}
		},
	)
	bot.Handle(
		"/filter", func(context telebot.Context) error {
			return sendFilter(context, namesFilter)
		},
	)
	bot.Handle(
		"/filter_add", func(context telebot.Context) error {
			for _, arg := range context.Args() {
				namesFilter = append(namesFilter, arg)
			}
			return sendFilter(context, namesFilter)
		},
	)
	bot.Handle(
		"/filter_remove", func(context telebot.Context) error {
			for _, arg := range context.Args() {
				index := slices.Index(namesFilter, arg)
				if index >= 0 {
					namesFilter = slices.Delete(namesFilter, index, index+1)
				}
			}
			return sendFilter(context, namesFilter)
		},
	)
	bot.Handle(
		"/processes", func(context telebot.Context) error {
			printProcessList(bot, admin, namesFilter)
			return nil
		},
	)
	err = bot.SetCommands(
		[]telebot.Command{
			{
				Text:        "start",
				Description: "Start the bot",
			},
			{
				Text:        "hello",
				Description: "Greeting for you",
			},
			{
				Text:        "filter",
				Description: "Print current filter",
			},
			{
				Text:        "filter_add",
				Description: "Add new value to filter",
			},
			{
				Text:        "filter_remove",
				Description: "Remove a value from filter",
			}, {
				Text:        "processes",
				Description: "Prints processes list",
			},
		},
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Start()
}

func sendFilter(context telebot.Context, namesFilter []string) error {
	message := strings.Builder{}
	if len(namesFilter) == 0 {
		return context.Send("Фильтр пуст")
	}

	message.WriteString("Текущий фильтр:\n")
	for _, filter := range namesFilter {
		message.WriteString(fmt.Sprintf("%v\n", filter))
	}
	return context.Send(message.String())
}

func printProcessList(bot *telebot.Bot, admin *telebot.User, filter []string) {
	processes, err := process.Processes()
	if err != nil {
		log.Fatal(err)
	}
	var allProcesses []ProcessInfo
	for _, p := range processes {
		processName, _ := p.Name()
		processExe, _ := p.Exe()
		terminal, _ := p.Terminal()
		if !slices.Contains(filter, strings.ToLower(processName)) {
			continue
		}

		allProcesses = append(
			allProcesses, ProcessInfo{
				Pid:      p.Pid,
				Name:     processName,
				Exe:      processExe,
				Terminal: terminal,
			},
		)
	}
	message := strings.Builder{}
	message.WriteString("Here are all the processes:\n")
	for _, p := range allProcesses {
		message.WriteString(fmt.Sprintf("%v %v - \"%v\"\n", p.Pid, p.Name, p.Exe))
	}
	_, _ = bot.Send(admin, message.String())
}

type ProcessInfo struct {
	Pid      int32
	Name     string
	Exe      string
	Terminal string
}
