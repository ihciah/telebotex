package auth

import (
	"fmt"
	"log"

	"github.com/ihciah/telebotex"
	"github.com/ihciah/telebotex/plugin"
	jsoniter "github.com/json-iterator/go"
	tb "gopkg.in/tucnak/telebot.v2"
)

const authenticatorKey = "authenticator"

type Authenticator struct {
	Admin []int `json:"admin"`
}

func NewAuthenticator(bot *telebotex.Bot) (*Authenticator, error) {
	a := new(Authenticator)
	err := bot.InitInterceptor(a)
	return a, err
}

func (a *Authenticator) LoadConfig(config map[string]jsoniter.RawMessage) error {
	err := telebotex.UnmarshalFromConfig(config, authenticatorKey, a)
	return err
}

func (a *Authenticator) Wrap(bot plugin.TelegramBot) plugin.TelegramBot {
	admins := make(map[int]struct{}, len(a.Admin))
	for _, id := range a.Admin {
		admins[id] = struct{}{}
	}
	return &botWrapped{
		TelegramBot: bot,
		admins:      admins,
	}
}

type botWrapped struct {
	plugin.TelegramBot
	admins map[int]struct{}
}

func (b *botWrapped) Handle(endpoint interface{}, handler interface{}) {
	if handlerFunc, ok := handler.(func(m *tb.Message)); ok {
		b.TelegramBot.Handle(endpoint, func(m *tb.Message) {
			if _, ok := b.admins[m.Sender.ID]; ok {
				handlerFunc(m)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", m.Sender.ID)
				log.Print(msg)
				_, _ = b.Send(m.Sender, msg)
			}
		})
	} else {
		log.Print("handler signature is not supported")
	}
}
