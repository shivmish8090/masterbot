package modules

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

const (
	MaxHandlers    = 20
	MaxHelpModules = 10
)

var (
	Continue    = ext.ContinueGroups
	Handlers    = make([]ext.Handler, 0, MaxHandlers)
	ModulesHelp = make(map[string]struct {
		Callback string
		Help     string
	}, MaxHelpModules)
)

func orCont(err error) any {
	if err != nil {
		return err
	}
	return Continue
}

func Register(h ext.Handler) {
	if len(Handlers) >= MaxHandlers {
		log.Panic("handler limit exceeded")
	}
	Handlers = append(Handlers, h)
}

func AddHelp(name, callback, help string, h ext.Handler) {
	if len(ModulesHelp) >= MaxHelpModules {
		log.Panic("help modules limit exceeded")
	}
	if h == nil {
		h = handlers.NewCallback(callbackquery.Equal(callback), helpModuleCB)
	}
	Register(h)
	ModulesHelp[name] = struct {
		Callback string
		Help     string
	}{
		Callback: callback,
		Help:     help,
	}
}

func GetHelp(callback string) string {
	for _, data := range ModulesHelp {
		if data.Callback == callback {
			return data.Help
		}
	}
	return ""
}
