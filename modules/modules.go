package modules

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

var (
	Handlers    []ext.Handler
	ModulesHelp = map[string]struct {
		Callback string
		Help     string
	}{}
)

func Register(h ext.Handler) {
	Handlers = append(Handlers, h)
}

func AddHelp(name, callback, help string, h ext.Handler) {
	ModulesHelp[name] = struct {
		Callback string
		Help     string
	}{
		Callback: callback,
		Help:     help,
	}

	var handler ext.Handler
	if h != nil {
		handler = h
	} else {
		handler = handlers.NewCallback(callbackquery.Equal(callback), helpModuleCB)
	}

	Register(handler)
}
