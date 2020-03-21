package interceptor

import (
	"log"

	"github.com/ihciah/telebotex/bot"
	"github.com/ihciah/telebotex/plugin"
	jsoniter "github.com/json-iterator/go"
)

type Interceptor interface {
	LoadConfig(map[string]jsoniter.RawMessage) error
	Wrap(bot.TelegramBotExt) bot.TelegramBotExt
}

type BaseInterceptor struct{}

func (b *BaseInterceptor) LoadConfig(_ map[string]jsoniter.RawMessage) error {
	log.Println("ignore json config")
	return nil
}

func (b *BaseInterceptor) Wrap(_ bot.TelegramBotExt) bot.TelegramBotExt {
	log.Panic("Wrap interface is not implemented")
	return nil
}

type InterceptedPlugin struct {
	plugin.Plugin
	Interceptors []Interceptor
}

func NewInterceptedPlugin(plugin plugin.Plugin, ic ...Interceptor) *InterceptedPlugin {
	return &InterceptedPlugin{
		Plugin:       plugin,
		Interceptors: ic,
	}
}
