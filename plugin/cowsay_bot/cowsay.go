package cowsay_bot

import (
	"fmt"
	"strings"

	"github.com/ihciah/telebotex/bot"
	"github.com/ihciah/telebotex/plugin"
	tb "gopkg.in/tucnak/telebot.v2"
)

const cow = `
\   ^__^
 \  (oo)\_______
    (__)\       )\/\
        ||----w |
        ||     ||`

const sheep = `
\    __     
 \  UooU\.'@@@@@@'.
    \__/(@@@@@@@@@@)
         (@@@@@@@@)
         'YY~~~~YY'
          ||    ||
`

type Bot struct {
	plugin.BasePlugin
}

func (b *Bot) Register(bot bot.TelegramBotExt) {
	bot.ConditionalHandle(tb.OnText, func(m *tb.Message) {
		_, _ = bot.Send(m.Sender, say(m.Text[6:], cow))
	}, isCowSay)
	bot.ConditionalHandle(tb.OnText, func(m *tb.Message) {
		_, _ = bot.Send(m.Sender, say(m.Text[8:], sheep))
	}, isSheepSay)
	bot.Handle(tb.OnText, func(m *tb.Message) {
		_, _ = bot.Send(m.Sender, "This is a conditional handler example, "+
			"please reply `cowsay $name` or `sheepsay $name`")
	})
}

func isCowSay(m *tb.Message) bool {
	return strings.HasPrefix(m.Text, "cowsay")
}

func isSheepSay(m *tb.Message) bool {
	return strings.HasPrefix(m.Text, "sheepsay")
}

func say(name, animal string) string {
	return fmt.Sprintf("Hello %s%s", name, animal)
}

func NewBot() *Bot {
	return new(Bot)
}
