package telebotex

import (
	"log"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type condHandler struct {
	Handler  interface{}
	CondFunc interface{}
}

type catchAllHandler struct {
	CondHandlers []condHandler
	CatchAll     interface{}
}

type botExt struct {
	*tb.Bot
	conditionalHandlers map[string]catchAllHandler
}

func newTelegramBotExt(config *telegramConfig) (*botExt, error) {
	if config.Token == "" {
		return nil, configLoadErr
	}

	tbBot, err := tb.NewBot(tb.Settings{
		Token:  config.Token,
		URL:    config.URL,
		Poller: &tb.LongPoller{Timeout: time.Duration(config.PollTimeout) * time.Second},
	})
	if err != nil {
		return nil, err
	}
	return &botExt{
		Bot:                 tbBot,
		conditionalHandlers: make(map[string]catchAllHandler),
	}, nil
}

func (b *botExt) ConditionalHandle(endpoint interface{}, handler interface{}, condFunc interface{}) {
	// condFunc must be some type, ref: telebot.go line 35
	if endpointStr, ok := endpoint.(string); ok && strings.HasPrefix(endpointStr, "\a") {
		catchAll := b.conditionalHandlers[endpointStr]
		catchAll.CondHandlers = append(catchAll.CondHandlers, condHandler{
			Handler:  handler,
			CondFunc: condFunc,
		})
		b.conditionalHandlers[endpointStr] = catchAll
	} else {
		log.Panic("conditional handler endpoint must be special string, like tb.OnText")
	}
}

func (b *botExt) Handle(endpoint interface{}, handler interface{}) {
	if endpointStr, ok := endpoint.(string); ok && strings.HasPrefix(endpointStr, "\a") {
		// it's a catch-all handler
		catchAll := b.conditionalHandlers[endpointStr]
		catchAll.CatchAll = handler
		b.conditionalHandlers[endpointStr] = catchAll
	} else {
		b.Bot.Handle(endpoint, handler)
	}
}

func (b *botExt) registerCatchAll() {
	// Register catch-all handlers
	for endpoint, catchAll := range b.conditionalHandlers {
		var catchAllFunc interface{}
		switch endpoint {
		case tb.OnText, tb.OnPhoto, tb.OnAudio, tb.OnDocument, tb.OnSticker, tb.OnVideo, tb.OnVoice, tb.OnVideoNote,
			tb.OnContact, tb.OnLocation, tb.OnVenue, tb.OnEdited, tb.OnPinned, tb.OnChannelPost, tb.OnEditedChannelPost,
			tb.OnAddedToGroup, tb.OnUserJoined, tb.OnUserLeft, tb.OnNewGroupTitle, tb.OnNewGroupPhoto, tb.OnGroupPhotoDeleted:
			catchAllFunc = registerCatchAllMessage(catchAll)
		case tb.OnMigration:
			catchAllFunc = registerCatchAllInt64(catchAll)
		case tb.OnCallback:
			catchAllFunc = registerCatchAllCallback(catchAll)
		case tb.OnQuery:
			catchAllFunc = registerCatchAllQuery(catchAll)
		case tb.OnChosenInlineResult:
			catchAllFunc = registerCatchAllChosenInlineResult(catchAll)
		case tb.OnCheckout:
			catchAllFunc = registerCatchAllPreCheckoutQuery(catchAll)
		}

		if catchAllFunc != nil {
			b.Bot.Handle(endpoint, catchAllFunc)
		}
	}
}

func registerCatchAllMessage(catchAll catchAllHandler) func(m *tb.Message) {
	type messageFuncPair struct {
		Cf func(*tb.Message) bool
		Ch func(*tb.Message)
	}
	pairs := make([]messageFuncPair, 0, len(catchAll.CondHandlers)+1)
	for _, condH := range catchAll.CondHandlers {
		cf, okCf := condH.CondFunc.(func(*tb.Message) bool)
		ch, okCh := condH.Handler.(func(*tb.Message))
		if okCf && okCh {
			pairs = append(pairs, messageFuncPair{
				Cf: cf,
				Ch: ch,
			})
		}
	}
	if ca, okCa := catchAll.CatchAll.(func(*tb.Message)); okCa && ca != nil {
		pairs = append(pairs, messageFuncPair{
			Cf: func(_ *tb.Message) bool {
				return true
			},
			Ch: ca,
		})
	}
	if len(pairs) != 0 {
		return func(m *tb.Message) {
			for _, p := range pairs {
				if p.Cf(m) {
					p.Ch(m)
					break
				}
			}
		}
	}
	return nil
}

func registerCatchAllInt64(catchAll catchAllHandler) func(from, to int64) {
	type messageFuncPair struct {
		Cf func(from, to int64) bool
		Ch func(from, to int64)
	}
	pairs := make([]messageFuncPair, 0, len(catchAll.CondHandlers)+1)
	for _, condH := range catchAll.CondHandlers {
		cf, okCf := condH.CondFunc.(func(from, to int64) bool)
		ch, okCh := condH.Handler.(func(from, to int64))
		if okCf && okCh {
			pairs = append(pairs, messageFuncPair{
				Cf: cf,
				Ch: ch,
			})
		}
	}
	if ca, okCa := catchAll.CatchAll.(func(from, to int64)); okCa && ca != nil {
		pairs = append(pairs, messageFuncPair{
			Cf: func(_, _ int64) bool {
				return true
			},
			Ch: ca,
		})
	}
	if len(pairs) != 0 {
		return func(from, to int64) {
			for _, p := range pairs {
				if p.Cf(from, to) {
					p.Ch(from, to)
					break
				}
			}
		}
	}
	return nil
}

func registerCatchAllCallback(catchAll catchAllHandler) func(c *tb.Callback) {
	type messageFuncPair struct {
		Cf func(c *tb.Callback) bool
		Ch func(c *tb.Callback)
	}
	pairs := make([]messageFuncPair, 0, len(catchAll.CondHandlers)+1)
	for _, condH := range catchAll.CondHandlers {
		cf, okCf := condH.CondFunc.(func(c *tb.Callback) bool)
		ch, okCh := condH.Handler.(func(c *tb.Callback))
		if okCf && okCh {
			pairs = append(pairs, messageFuncPair{
				Cf: cf,
				Ch: ch,
			})
		}
	}
	if ca, okCa := catchAll.CatchAll.(func(c *tb.Callback)); okCa && ca != nil {
		pairs = append(pairs, messageFuncPair{
			Cf: func(_ *tb.Callback) bool {
				return true
			},
			Ch: ca,
		})
	}
	if len(pairs) != 0 {
		return func(c *tb.Callback) {
			for _, p := range pairs {
				if p.Cf(c) {
					p.Ch(c)
					break
				}
			}
		}
	}
	return nil
}

func registerCatchAllQuery(catchAll catchAllHandler) func(q *tb.Query) {
	type messageFuncPair struct {
		Cf func(q *tb.Query) bool
		Ch func(q *tb.Query)
	}
	pairs := make([]messageFuncPair, 0, len(catchAll.CondHandlers)+1)
	for _, condH := range catchAll.CondHandlers {
		cf, okCf := condH.CondFunc.(func(q *tb.Query) bool)
		ch, okCh := condH.Handler.(func(q *tb.Query))
		if okCf && okCh {
			pairs = append(pairs, messageFuncPair{
				Cf: cf,
				Ch: ch,
			})
		}
	}
	if ca, okCa := catchAll.CatchAll.(func(q *tb.Query)); okCa && ca != nil {
		pairs = append(pairs, messageFuncPair{
			Cf: func(_ *tb.Query) bool {
				return true
			},
			Ch: ca,
		})
	}
	if len(pairs) != 0 {
		return func(q *tb.Query) {
			for _, p := range pairs {
				if p.Cf(q) {
					p.Ch(q)
					break
				}
			}
		}
	}
	return nil
}

func registerCatchAllChosenInlineResult(catchAll catchAllHandler) func(c *tb.ChosenInlineResult) {
	type messageFuncPair struct {
		Cf func(c *tb.ChosenInlineResult) bool
		Ch func(c *tb.ChosenInlineResult)
	}
	pairs := make([]messageFuncPair, 0, len(catchAll.CondHandlers)+1)
	for _, condH := range catchAll.CondHandlers {
		cf, okCf := condH.CondFunc.(func(c *tb.ChosenInlineResult) bool)
		ch, okCh := condH.Handler.(func(c *tb.ChosenInlineResult))
		if okCf && okCh {
			pairs = append(pairs, messageFuncPair{
				Cf: cf,
				Ch: ch,
			})
		}
	}
	if ca, okCa := catchAll.CatchAll.(func(c *tb.ChosenInlineResult)); okCa && ca != nil {
		pairs = append(pairs, messageFuncPair{
			Cf: func(_ *tb.ChosenInlineResult) bool {
				return true
			},
			Ch: ca,
		})
	}
	if len(pairs) != 0 {
		return func(c *tb.ChosenInlineResult) {
			for _, p := range pairs {
				if p.Cf(c) {
					p.Ch(c)
					break
				}
			}
		}
	}
	return nil
}

func registerCatchAllPreCheckoutQuery(catchAll catchAllHandler) func(p *tb.PreCheckoutQuery) {
	type messageFuncPair struct {
		Cf func(p *tb.PreCheckoutQuery) bool
		Ch func(p *tb.PreCheckoutQuery)
	}
	pairs := make([]messageFuncPair, 0, len(catchAll.CondHandlers)+1)
	for _, condH := range catchAll.CondHandlers {
		cf, okCf := condH.CondFunc.(func(p *tb.PreCheckoutQuery) bool)
		ch, okCh := condH.Handler.(func(p *tb.PreCheckoutQuery))
		if okCf && okCh {
			pairs = append(pairs, messageFuncPair{
				Cf: cf,
				Ch: ch,
			})
		}
	}
	if ca, okCa := catchAll.CatchAll.(func(p *tb.PreCheckoutQuery)); okCa && ca != nil {
		pairs = append(pairs, messageFuncPair{
			Cf: func(_ *tb.PreCheckoutQuery) bool {
				return true
			},
			Ch: ca,
		})
	}
	if len(pairs) != 0 {
		return func(pr *tb.PreCheckoutQuery) {
			for _, p := range pairs {
				if p.Cf(pr) {
					p.Ch(pr)
					break
				}
			}
		}
	}
	return nil
}

func (b *botExt) Start() {
	b.registerCatchAll()
	b.Bot.Start()
}
