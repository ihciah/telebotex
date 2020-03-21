package telebotex

import (
	"log"
	"time"

	"github.com/ihciah/telebotex/bot"
	"github.com/ihciah/telebotex/interceptor"
	"github.com/ihciah/telebotex/plugin"
	jsoniter "github.com/json-iterator/go"
)

type Bot struct {
	bot.TelegramBotExt
	configs map[string]jsoniter.RawMessage
}

func NewBot(configFile string, retry bool) (*Bot, error) {
	config, err := newConfig(configFile)
	if err != nil {
		log.Printf("config loading error: %v", err)
		return nil, err
	}
	log.Print("config loading successfully")

	telegramConfig, err := newTelegramConfig(config)
	if err != nil {
		log.Printf("telegram config loading error: %v", err)
		return nil, err
	}
	log.Print("telegram config loading successfully")

	var b *botExt
	for {
		b, err = newTelegramBotExt(telegramConfig)
		if !retry || err == nil {
			log.Print("telegram bot creating successfully")
			break
		}
		log.Printf("telegram bot creating error: %v", err)
		log.Print("will sleep and retry")
		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, err
	}
	return &Bot{
		TelegramBotExt: b,
		configs:        config,
	}, nil
}

func MustNewBot(configFile string) *Bot {
	b, err := NewBot(configFile, true)
	if err != nil {
		panic(err)
	}
	return b
}

func (b *Bot) Register(plugins ...plugin.Plugin) error {
	for _, p := range plugins {
		if err := p.LoadConfig(b.configs); err != nil {
			log.Printf("plugin config loading error: %v", err)
			return err
		} else {
			log.Print("plugin config loading successfully")
		}
		if intercepted, ok := p.(*interceptor.InterceptedPlugin); ok {
			wrap := b.TelegramBotExt
			for _, i := range intercepted.Interceptors {
				wrap = i.Wrap(wrap)
			}
			p.Register(wrap)
		} else {
			p.Register(b)
		}
	}
	return nil
}

func (b *Bot) InitInterceptor(interceptors ...interceptor.Interceptor) error {
	for _, inte := range interceptors {
		err := inte.LoadConfig(b.configs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) Start() {
	log.Print("bot start")
	b.TelegramBotExt.Start()
}
