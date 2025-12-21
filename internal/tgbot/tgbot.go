package tgbot

import (
	"context"
	"errors"
	"fmt"
	"game-controller-bot/internal/botcommand"
	"game-controller-bot/internal/botcontext"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/creativeprojects/go-selfupdate"
	"gopkg.in/telebot.v4"
)

var ErrWrongBotToken = errors.New("wrong telegram bot token")
var ErrFailedToSetCommands = errors.New("failed to set up bot commands")
var ErrFailedToSelfUpdate = errors.New("failed to self update bot")

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
			log.Println(fmt.Errorf("failed to handle: %w", err))
			if context != nil {
				err2 := context.Send(fmt.Sprintf("Ошибка: %v", err.Error()))
				if err2 != nil {
					log.Println("Failed to send error response", err2)
				}
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

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		b.bot.Stop()
	}()

	b.bot.Start()
	return nil
}

func (b *TgBot) StartSelfUpdate(version string) {
	if version == "0.0.0" {
		log.Println("Skip self update for developer version")
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := b.selfUpdate(version)
				if err != nil {
					log.Println(fmt.Errorf("%w: %w", ErrFailedToSelfUpdate, err))
				}
			case <-quit:
				return
			}
		}
	}()
}

func (b *TgBot) selfUpdate(version string) error {
	latest, found, err := selfupdate.DetectLatest(
		context.Background(), selfupdate.ParseSlug("mih-kopylov/game-controller-bot"),
	)
	if err != nil {
		return fmt.Errorf("failed to detect latest version: %w", err)
	}
	if !found {
		return fmt.Errorf("didn't found latest version for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	if latest.LessOrEqual(version) {
		log.Printf("Current version %s is the latest\n", version)
		return nil
	}

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("could not locate executable path: %w", err)
	}
	err = selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe)
	if err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start new binary version: %w", err)
	}
	log.Printf("Successfully updated to version %s, so exit\n", latest.Version())
	os.Exit(0)
	return nil
}
