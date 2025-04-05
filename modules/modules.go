package modules

import "github.com/PaulSonOfLars/gotgbot/v2/ext"

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

func AddHelp(name, callback, help string) {
	ModulesHelp[name] = struct {
		Callback string
		Help     string
	}{
		Callback: callback,
		Help:     help,
	}
}
