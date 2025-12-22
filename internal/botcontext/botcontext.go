package botcontext

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var ErrFailedToReadContextData = errors.New("failed to read context data")
var ErrFailedToWriteContextData = errors.New("failed to write context data")

const contextDataFileName = "game-controller-bot.yaml"

type BotContext struct {
}

func NewBotContext() *BotContext {
	return &BotContext{}
}

func (b *BotContext) Read() (*Data, error) {
	path, err := b.getFullFilePath()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToReadContextData, err)
	}

	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Data{
			Admins:      nil,
			SystemUser:  "",
			NamesToKill: nil,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToReadContextData, err)
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToReadContextData, err)
	}

	var result Data
	err = yaml.Unmarshal(bytes, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToReadContextData, err)
	}

	return &result, nil
}

func (b *BotContext) Write(data *Data) error {
	path, err := b.getFullFilePath()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToWriteContextData, err)
	}

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToWriteContextData, err)
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToWriteContextData, err)
	}

	return nil
}

func (b *BotContext) getFullFilePath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	fullFilePath := filepath.Join(userConfigDir, contextDataFileName)
	return fullFilePath, nil
}
