package internal

import (
	"errors"
	"fmt"
	"game-controller-bot/internal/botcommand"
	"game-controller-bot/internal/tgbot"
)

var ErrFailedToStartBot = errors.New("failed to start bot")

func StartBot() error {
	bot, err := tgbot.NewTgBot()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToStartBot, err)
	}

	bot.AddCommand(botcommand.NewStartCommand(bot.Context))

	bot.AddCommand(botcommand.NewUserListCommand(bot.Context))
	bot.AddCommand(botcommand.NewUserSetCommand(bot.Context))
	bot.AddCommand(botcommand.NewUserClearCommand(bot.Context))

	bot.AddCommand(botcommand.NewFilterListCommand(bot.Context))
	bot.AddCommand(botcommand.NewFilterAddCommand(bot.Context))
	bot.AddCommand(botcommand.NewFilterRemoveCommand(bot.Context))

	bot.AddCommand(botcommand.NewProcessListCommand(bot.Context))
	bot.AddCommand(botcommand.NewProcessTerminateCommand(bot.Context))

	err = bot.Start()
	if err != nil {
		return err
	}

	return nil
}
