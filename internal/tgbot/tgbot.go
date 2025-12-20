package tgbot

import (
	"errors"
	"fmt"
	"game-controller-bot/internal/botcommand"
	"game-controller-bot/internal/botcontext"
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v4"
)

var ErrWrongBotToken = errors.New("wrong telegram bot token")
var ErrFailedToSetCommands = errors.New("failed to set up bot commands")

type TgBot struct {
	bot      *telebot.Bot
	commands []telebot.Command
	Context  *botcontext.BotContext
}

func NewTgBot() (*TgBot, error) {
	prefs := telebot.Settings{
		Token:  os.Getenv("CONTROLLER_BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, context telebot.Context) {
			if context != nil {
				err2 := context.Send(err.Error())
				if err2 != nil {
					log.Println("Failed to send error response", err2)
					log.Println("Cause error", err)
				}
			} else {
				log.Println("Failed to handle", err)
			}
		},
	}

	bot, err := telebot.NewBot(prefs)
	if err != nil {
		if errors.Is(err, telebot.ErrNotFound) {
			return nil, fmt.Errorf("%w: %w", ErrWrongBotToken, err)
		}
		return nil, err
	}

	result := &TgBot{
		bot:     bot,
		Context: botcontext.NewBotContext(),
	}

	return result, nil
}

func (b *TgBot) AddCommand(botCommand botcommand.BotCommand) {
	b.bot.Handle(botCommand.GetName(), botCommand.GetHandleFunc())
	b.commands = append(
		b.commands, telebot.Command{
			Text:        botCommand.GetName(),
			Description: botCommand.GetDescription(),
		},
	)
}

func (b *TgBot) Start() error {
	err := b.bot.SetCommands(b.commands)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSetCommands, err)
	}

	b.bot.Start()
	return nil
}
