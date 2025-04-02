package main

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

func EvalHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if len(ctx.Args()) < 2 {
		ctx.EffectiveMessage.Reply(b, "Usage: /eval <go code>", nil)
		return nil
	}

	code := strings.SplitN(ctx.EffectiveMessage.GetText(), " ", 2)[1]
	cleanCode, imports := extractImportsAndCode(code)

	/*botString, err := json.Marshal(b)
	if err != nil {
		ctx.EffectiveMessage.Reply(b, "Error: Failed to serialize bot", nil)
		return nil
	}*/

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

func extractImportsAndCode(code string) (string, string) {
	var imports []string
	importRegex := regexp.MustCompile(`(?m)^\s*import\s+(\"[^"]+\"|[a-zA-Z0-9_]+?\s+"[^"]+")`)

	matches := importRegex.FindAllString(code, -1)
	for _, match := range matches {
		imports = append(imports, strings.TrimSpace(match))
	}

	cleanCode := importRegex.ReplaceAllString(code, "")
	formattedImports := strings.Join(imports, "\n")

	return strings.TrimSpace(cleanCode), formattedImports
}

func runGoCode(code, imports, botString, ctxString string) (string, error) {
	evalTemplate := `package main

import (
	"encoding/json"
	"fmt"
	%s
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
        "github.com/Vivekkumar-IN/EditguardianBot/config"
)

var ctxString = %q

func main() {
	var ctx ext.Context
 
        Bot, err := gotgbot.NewBot(config.Token,  nil)
        if err != nil {
                panic("failed to create new bot: " + err.Error())
        }

	json.Unmarshal([]byte(ctxString), &ctx)

	%s
}
`
	evalCode := fmt.Sprintf(evalTemplate, imports, ctxString, code)

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
