package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)
// Eval code

func EvalHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if len(ctx.Args()) < 2 {
		_, _ = ctx.EffectiveMessage.Reply(b, "Usage: /eval <go code>", nil)
		return nil
	}

	code := strings.Join(ctx.Args()[1:], "\n")
	cleanCode, imports := extractImportsAndCode(code)

	result, err := runGoCode(cleanCode, imports, b, ctx)
	if err != nil {
		result = "Error: " + err.Error()
	}

	_, _ = ctx.EffectiveMessage.Reply(b, result, nil)
	return nil
}

func extractImportsAndCode(code string) (string, string) {
	var imports []string
	importRegex := regexp.MustCompile(`(?m)^\s*import\s+(.+?|"[^"]+"|[a-zA-Z0-9_]+?\s+"[^"]+")`)

	matches := importRegex.FindAllString(code, -1)
	for _, match := range matches {
		imports = append(imports, strings.TrimSpace(match))
	}

	cleanCode := importRegex.ReplaceAllString(code, "")
	formattedImports := strings.Join(imports, "\n")

	return strings.TrimSpace(cleanCode), formattedImports
}

func runGoCode(code, imports string, b *gotgbot.Bot, ctx *ext.Context) (string, error) {
	evalTemplate := `
		package evalpkg

		import (
			"fmt"
			"github.com/PaulSonOfLars/gotgbot/v2"
			"github.com/PaulSonOfLars/gotgbot/v2/ext"
			%s
		)

		func EvalCode(b *gotgbot.Bot, ctx *ext.Context) string {
			var output bytes.Buffer
			fmt.Fprintln(&output, func() string {
				%s
			}())
			return output.String()
		}
	`

	evalCode := fmt.Sprintf(evalTemplate, imports, code)

	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	_, err := i.Eval(evalCode)
	if err != nil {
		return "", err
	}

	v, err := i.Eval("evalpkg.EvalCode")
	if err != nil {
		return "", err
	}

	evalFunc := v.Interface().(func(*gotgbot.Bot, *ext.Context) string)
	return evalFunc(b, ctx), nil
}