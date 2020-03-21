package telebotex

import (
	jsoniter "github.com/json-iterator/go"
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
