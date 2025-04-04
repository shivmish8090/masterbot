package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func init(){

        RegisterCommand(handlers.NewMessage(
                filters.AndFilter(filters.Owner, filters.Command(b, "eval")),
                EvalHandler,
        ).SetAllowEdited(true))
}


func EvalHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if len(ctx.Args()) < 2 {
		ctx.EffectiveMessage.Reply(b, "Usage: /eval <go code>", nil)
		return nil
	}

	code := strings.SplitN(ctx.EffectiveMessage.GetText(), " ", 2)[1]
	cleanCode, imports := resolveImports(code)

	ctxString, err := json.Marshal(ctx)
	if err != nil {
		ctx.EffectiveMessage.Reply(b, "Error: Failed to serialize context "+err.Error(), nil)
		return nil
	}

	result, err := runGoCode(cleanCode, imports, string(ctxString))
	if err != nil {
		result = "Error: " + err.Error()
	}

	ctx.EffectiveMessage.Reply(b, result, nil)
	return nil
}

func resolveImports(code string) (string, []string) {
	var imports []string
	importsRegex := regexp.MustCompile(`import\s*([\s\S]*?)|import\s*"([^"]+)"`)
	importsMatches := importsRegex.FindAllStringSubmatch(code, -1)

	for _, v := range importsMatches {
		if v[1] != "" { // Multi-line import
			lines := strings.Split(v[1], "\n")
			for _, line := range lines {
				cleanLine := strings.TrimSpace(line)
				if cleanLine != "" {
					imports = append(imports, fmt.Sprintf("%q", cleanLine))
				}
			}
		} else if v[2] != "" { // Single-line import
			imports = append(imports, fmt.Sprintf("%q", v[2]))
		}
	}

	code = importsRegex.ReplaceAllString(code, "")
	return code, imports
}

func runGoCode(code string, imports []string, ctxString string) (string, error) {
	var importBlock string
	if len(imports) > 0 {
		importBlock = fmt.Sprintf(`import (
    "encoding/json"
    "fmt"
    %s
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/Vivekkumar-IN/EditguardianBot/config"
)`, strings.Join(imports, "\n    "))
	} else {
		importBlock = `import (
    "encoding/json"
    "fmt"
    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
    "github.com/Vivekkumar-IN/EditguardianBot/config"
)`
	}

	evalTemplate := `package main

%s

var ctxString = %q

func main() {
    var ctx ext.Context

    Bot, err := gotgbot.NewBot(config.Token, nil)
    if err != nil {
        panic("failed to create new bot: " + err.Error())
    }

    json.Unmarshal([]byte(ctxString), &ctx)
    b := Bot
    print := fmt.Println
    printf := fmt.Printf

    %s

    _ = ctx
    _ = b
    _ = print
    _ = printf
}`
	evalCode := fmt.Sprintf(evalTemplate, importBlock, ctxString, code)

	tmpFile := fmt.Sprintf("/tmp/eval_%d.go", time.Now().UnixNano())
	err := os.WriteFile(tmpFile, []byte(evalCode), 0o644)
	if err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tmpFile)

	cmd := exec.Command("go", "run", tmpFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %w", string(output), err)
	}

	return string(output), nil
}
