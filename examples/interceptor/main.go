package main

import (
	"github.com/ihciah/telebotex"
	"github.com/ihciah/telebotex/interceptor"
	"github.com/ihciah/telebotex/interceptor/auth"
	"github.com/ihciah/telebotex/interceptor/private"
	"github.com/ihciah/telebotex/plugin/id_bot"
)

func main() {
	bot := telebotex.MustNewBot("config.json")
	authInterceptor, err := auth.NewAuthenticator(bot)
	privateInterceptor := private.NewPrivateChatFilter()
	if err != nil {
		panic(err)
	}
	plugin := interceptor.NewInterceptedPlugin(id_bot.NewBot(), privateInterceptor, authInterceptor)

	if err := bot.Register(plugin); err != nil {
		panic(err)
	}

	bot.Start()
}
