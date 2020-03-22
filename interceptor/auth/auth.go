/*
Authenticator interceptor hooks Handle and ConditionalHandle function.
It will try to get user id and check if it is in admin list.
Attention: this module only works for part of handler functions:
	func(m *tb.Message), func(c *tb.Callback), func(q *tb.Query), func(c *tb.ChosenInlineResult)
*/
package auth

import (
	"fmt"
	"log"

	"github.com/ihciah/telebotex"
	"github.com/ihciah/telebotex/bot"
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

func (a *Authenticator) Wrap(bot bot.TelegramBotExt) bot.TelegramBotExt {
	admins := make(map[int]struct{}, len(a.Admin))
	for _, id := range a.Admin {
		admins[id] = struct{}{}
	}
	return &botWrapped{
		TelegramBotExt: bot,
		admins:         admins,
	}
}

type botWrapped struct {
	bot.TelegramBotExt
	admins map[int]struct{}
}

func (b *botWrapped) Handle(endpoint interface{}, handler interface{}) {
	switch h := handler.(type) {
	case func(m *tb.Message):
		b.TelegramBotExt.Handle(endpoint, func(m *tb.Message) {
			if _, ok := b.admins[m.Sender.ID]; ok {
				h(m)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", m.Sender.ID)
				log.Print(msg)
				_, _ = b.Send(m.Sender, msg)
			}
		})
	case func(c *tb.Callback):
		b.TelegramBotExt.Handle(endpoint, func(c *tb.Callback) {
			if _, ok := b.admins[c.Sender.ID]; ok {
				h(c)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", c.Sender.ID)
				log.Print(msg)
				_, _ = b.Send(c.Sender, msg)
			}
		})
	case func(q *tb.Query):
		b.TelegramBotExt.Handle(endpoint, func(q *tb.Query) {
			if _, ok := b.admins[q.From.ID]; ok {
				h(q)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", q.From.ID)
				log.Print(msg)
				_, _ = b.Send(&q.From, msg)
			}
		})
	case func(c *tb.ChosenInlineResult):
		b.TelegramBotExt.Handle(endpoint, func(c *tb.ChosenInlineResult) {
			if _, ok := b.admins[c.From.ID]; ok {
				h(c)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", c.From.ID)
				log.Print(msg)
				_, _ = b.Send(&c.From, msg)
			}
		})
	default:
		log.Print("handler signature is not supported")
	}
}

func (b *botWrapped) ConditionalHandle(endpoint interface{}, handler interface{}, condFunc interface{}) {
	switch h := handler.(type) {
	case func(m *tb.Message):
		b.TelegramBotExt.ConditionalHandle(endpoint, func(m *tb.Message) {
			if _, ok := b.admins[m.Sender.ID]; ok {
				h(m)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", m.Sender.ID)
				log.Print(msg)
				_, _ = b.Send(m.Sender, msg)
			}
		}, condFunc)
	case func(c *tb.Callback):
		b.TelegramBotExt.ConditionalHandle(endpoint, func(c *tb.Callback) {
			if _, ok := b.admins[c.Sender.ID]; ok {
				h(c)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", c.Sender.ID)
				log.Print(msg)
				_, _ = b.Send(c.Sender, msg)
			}
		}, condFunc)
	case func(q *tb.Query):
		b.TelegramBotExt.ConditionalHandle(endpoint, func(q *tb.Query) {
			if _, ok := b.admins[q.From.ID]; ok {
				h(q)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", q.From.ID)
				log.Print(msg)
				_, _ = b.Send(&q.From, msg)
			}
		}, condFunc)
	case func(c *tb.ChosenInlineResult):
		b.TelegramBotExt.ConditionalHandle(endpoint, func(c *tb.ChosenInlineResult) {
			if _, ok := b.admins[c.From.ID]; ok {
				h(c)
			} else {
				msg := fmt.Sprintf("Unauthorized user %d", c.From.ID)
				log.Print(msg)
				_, _ = b.Send(&c.From, msg)
			}
		}, condFunc)
	default:
		log.Print("handler signature is not supported")
	}
}
