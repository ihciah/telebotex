package interceptor

import (
	"log"

	"github.com/ihciah/telebotex/plugin"
	jsoniter "github.com/json-iterator/go"
)

type Interceptor interface {
	LoadConfig(map[string]jsoniter.RawMessage) error
	Wrap(plugin.TelegramBot) plugin.TelegramBot
}

type BaseInterceptor struct{}

func (b *BaseInterceptor) LoadConfig(_ map[string]jsoniter.RawMessage) error {
	log.Println("ignore json config")
	return nil
}

func (b *BaseInterceptor) Wrap(plugin.TelegramBot) plugin.TelegramBot {
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
