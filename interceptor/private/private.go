/*
PrivateChatFilter interceptor will filter out non-private chat
*/

package private

import (
	"fmt"
	"log"

	"github.com/ihciah/telebotex/bot"
	"github.com/ihciah/telebotex/interceptor"
	tb "gopkg.in/tucnak/telebot.v2"
)

type PrivateChatFilter struct {
	interceptor.BaseInterceptor
}

func NewPrivateChatFilter() *PrivateChatFilter {
	return new(PrivateChatFilter)
}

func (a *PrivateChatFilter) Wrap(bot bot.TelegramBotExt) bot.TelegramBotExt {
	return &botWrapped{bot}
}

type botWrapped struct {
	bot.TelegramBotExt
}

func (b *botWrapped) Handle(endpoint interface{}, handler interface{}) {
	if handlerFunc, ok := handler.(func(m *tb.Message)); ok {
		b.TelegramBotExt.Handle(endpoint, func(m *tb.Message) {
			if m.Private() {
				handlerFunc(m)
			} else {
				msg := fmt.Sprintf("Non-private chat %d", m.Sender.ID)
				log.Print(msg)
			}
		})
	} else {
		log.Print("handler signature is not supported")
	}
}

func (b *botWrapped) ConditionalHandle(endpoint interface{}, handler interface{}, condFunc interface{}) {
	if handlerFunc, ok := handler.(func(m *tb.Message)); ok {
		b.TelegramBotExt.ConditionalHandle(endpoint, func(m *tb.Message) {
			if m.Private() {
				handlerFunc(m)
			} else {
				msg := fmt.Sprintf("Non-private chat %d", m.Sender.ID)
				log.Print(msg)
			}
		}, condFunc)
	} else {
		log.Print("handler signature is not supported")
	}
}
