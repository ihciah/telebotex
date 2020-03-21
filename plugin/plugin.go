package plugin

import (
	"log"

	"github.com/ihciah/telebotex/bot"
	jsoniter "github.com/json-iterator/go"
)

type Plugin interface {
	LoadConfig(map[string]jsoniter.RawMessage) error
	Register(bot.TelegramBotExt)
}

type BasePlugin struct{}

func (p *BasePlugin) LoadConfig(_ map[string]jsoniter.RawMessage) error {
	log.Println("ignore json config")
	return nil
}

func (p *BasePlugin) Register(_ bot.TelegramBotExt) {
	log.Panic("Register interface is not implemented")
}
