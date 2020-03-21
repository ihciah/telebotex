package telebotex

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	tb "gopkg.in/tucnak/telebot.v2"
)

const telegramKey = "telegram"

type telegramConfig struct {
	Token       string `json:"token"`
	URL         string `json:"url"`
	PollTimeout int    `json:"poll_timeout"`
}

func newTelegramConfig(config map[string]jsoniter.RawMessage) (*telegramConfig, error) {
	cfg := new(telegramConfig)
	err := UnmarshalFromConfig(config, telegramKey, cfg)
	return cfg, err
}

func newTelegramBot(config *telegramConfig) (*tb.Bot, error) {
	if config.Token == "" {
		return nil, configLoadErr
	}

	return tb.NewBot(tb.Settings{
		Token:  config.Token,
		URL:    config.URL,
		Poller: &tb.LongPoller{Timeout: time.Duration(config.PollTimeout) * time.Second},
	})
}
