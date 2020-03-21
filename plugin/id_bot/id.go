package id_bot

import (
	"strconv"

	"github.com/ihciah/telebotex/bot"
	"github.com/ihciah/telebotex/plugin"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	plugin.BasePlugin
}

func (b *Bot) Register(bot bot.TelegramBotExt) {
	bot.Handle("/id", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		_, _ = bot.Send(m.Sender, strconv.Itoa(m.Sender.ID))
	})
}

func NewBot() *Bot {
	return new(Bot)
}
