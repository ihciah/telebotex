package main

import (
	"github.com/ihciah/telebotex"
	"github.com/ihciah/telebotex/plugin"
	"github.com/ihciah/telebotex/plugin/cowsay_bot"
)

func main() {
	bot := telebotex.MustNewBot("config.json")
	plugins := []plugin.Plugin{
		cowsay_bot.NewBot(),
	}

	err := bot.Register(plugins...)
	if err != nil {
		panic(err)
	}

	bot.Start()
}
