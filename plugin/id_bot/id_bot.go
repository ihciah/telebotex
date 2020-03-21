package id_bot

import (
	"strconv"

	"github.com/ihciah/telebotex/plugin"
	tb "gopkg.in/tucnak/telebot.v2"
)

type IdBot struct {
	plugin.BasePlugin
}

func (b *IdBot) Register(bot plugin.TelegramBot) {
	bot.Handle("/id", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		_, _ = bot.Send(m.Sender, strconv.Itoa(m.Sender.ID))
	})
}

func NewIdBot() *IdBot {
	return new(IdBot)
}
