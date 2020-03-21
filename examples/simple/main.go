package main

import (
	"github.com/ihciah/telebotex"
	"github.com/ihciah/telebotex/plugin"
	"github.com/ihciah/telebotex/plugin/id_bot"
)

func main() {
	bot := telebotex.MustNewBot("config.json")
	plugins := []plugin.Plugin{
		id_bot.NewBot(),
	}

	err := bot.Register(plugins...)
	if err != nil {
		panic(err)
	}

	bot.Start()
}
