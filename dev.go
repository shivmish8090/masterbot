package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)


func LsHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	text := ctx.EffectiveMessage.GetText()
	fields := strings.Fields(text)
	dir := "."
	if len(fields) > 1 {
		dir = strings.TrimSpace(strings.Replace(text, fields[0], "", 1))
	}

	cmd := exec.Command("ls", "-A", dir)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		ctx.EffectiveMessage.Reply(b, fmt.Sprintf("<b>Error:</b> <code>%s</code>", err.Error()), &gotgbot.SendMessageOpts{ParseMode: "HTML"})
		return nil
	}

	files := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(files) == 0 {
		ctx.EffectiveMessage.Reply(b, "<b>No files found in this directory.</b>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
		return nil
	}

	var responseBuilder strings.Builder
	var totalSize int64

	fileTypeEmoji := map[string]string{
		"file":   "📄",
		"dir":    "📂",
		"video":  "🎥",
		"audio":  "🎵",
		"image":  "🖼️",
		"go":     "🐹",
		"python": "🐍",
		"txt":    "📜",
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		fileType := "file"
		if fileInfo.IsDir() {
			fileType = "dir"
		} else {
			ext := strings.ToLower(filepath.Ext(file))
			switch ext {
			case ".mp4", ".mkv", ".webm", ".avi", ".flv", ".mov", ".wmv", ".3gp":
				fileType = "video"
			case ".mp3", ".wav", ".flac", ".ogg", ".m4a", ".wma":
				fileType = "audio"
			case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff":
				fileType = "image"
			case ".go":
				fileType = "go"
			case ".py":
				fileType = "python"
			case ".txt":
				fileType = "txt"
			}
		}

		fileSize := calcFileOrDirSize(filePath)
		totalSize += fileSize
		responseBuilder.WriteString(fmt.Sprintf("%s <b>%s</b> (%s)\n", fileTypeEmoji[fileType], file, sizeToHuman(fileSize)))
	}

	responseBuilder.WriteString(fmt.Sprintf("\n<b>Total Size:</b> %s", sizeToHuman(totalSize)))
	ctx.EffectiveMessage.Reply(b, responseBuilder.String(), &gotgbot.SendMessageOpts{ParseMode: "HTML"})

	return nil
}

func sizeToHuman(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d B", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	default:
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}

func calcFileOrDirSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}

	if !fi.IsDir() {
		return fi.Size()
	}

	var totalSize int64
	err = filepath.WalkDir(path, func(_ string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fi, err := info.Info()
			if err != nil {
				return err
			}
			totalSize += fi.Size()
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return totalSize
}

// Eval code

const boiler_code_for_eval = `
package main

import "fmt"
import "github.com/amarnathcjd/gogram/telegram"
import "encoding/json"

%s

var msg_id int32 = %d

var client *telegram.Client
var message *telegram.NewMessage
var m *telegram.NewMessage
var r *telegram.NewMessage
` + "var msg = `%s`\nvar snd = `%s`\nvar cht = `%s`\nvar chn = `%s`\nvar cch = `%s`" + `


func evalCode() {
        %s
}

func main() {
        var msg_o *telegram.MessageObj
        var snd_o *telegram.UserObj
        var cht_o *telegram.ChatObj
        var chn_o *telegram.Channel
        json.Unmarshal([]byte(msg), &msg_o)
        json.Unmarshal([]byte(snd), &snd_o)
        json.Unmarshal([]byte(cht), &cht_o)
        json.Unmarshal([]byte(chn), &chn_o)
        client, _ = telegram.NewClient(telegram.ClientConfig{
                StringSession: "%s",
        })

        client.Cache.ImportJSON([]byte(cch))

        client.Conn()

        x := []telegram.User{}
        y := []telegram.Chat{}
        x = append(x, snd_o)
        if chn_o != nil {
                y = append(y, chn_o)
        }
        if cht_o != nil {
                y = append(y, cht_o)
        }
        client.Cache.UpdatePeersToCache(x, y)
        idx := 0
        if cht_o != nil {
                idx = int(cht_o.ID)
        }
        if chn_o != nil {
                idx = int(chn_o.ID)
        }
        if snd_o != nil && idx == 0 {
                idx = int(snd_o.ID)
        }

        messageX, err := client.GetMessages(idx, &telegram.SearchOption{
                IDs: int(msg_id),
        })

        if err != nil {
                fmt.Println(err)
        }

        message = &messageX[0]
        m = message
        r, _ = message.GetReplyMessage()

        fmt.Println("output-start")
        evalCode()
}

func packMessage(c *telegram.Client, message telegram.Message, sender *telegram.UserObj, channel *telegram.Channel, chat *telegram.ChatObj) *telegram.NewMessage {
        var (
                m = &telegram.NewMessage{}
        )
        switch message := message.(type) {
        case *telegram.MessageObj:
                m.ID = message.ID
                m.OriginalUpdate = message
                m.Message = message
                m.Client = c
        default:
                return nil
        }
        m.Sender = sender
        m.Chat = chat
        m.Channel = channel
        if m.Channel != nil && (m.Sender.ID == m.Channel.ID) {
                m.SenderChat = channel
        } else {
                m.SenderChat = &telegram.Channel{}
        }
        m.Peer, _ = c.GetSendablePeer(message.(*telegram.MessageObj).PeerID)

        /*if m.IsMedia() {
                FileID := telegram.PackBotFileID(m.Media())
                m.File = &telegram.CustomFile{
                        FileID: FileID,
                        Name:   getFileName(m.Media()),
                        Size:   getFileSize(m.Media()),
                        Ext:    getFileExt(m.Media()),
                }
        }*/
        return m
}
`

func resolveImports(code string) (string, []string) {
	var imports []string
	importsRegex := regexp.MustCompile(`import\s*\(([\s\S]*?)\)|import\s*\"([\s\S]*?)\"`)
	importsMatches := importsRegex.FindAllStringSubmatch(code, -1)
	for _, v := range importsMatches {
		if v[1] != "" {
			imports = append(imports, v[1])
		} else {
			imports = append(imports, v[2])
		}
	}
	code = importsRegex.ReplaceAllString(code, "")
	return code, imports
}

func EvalHandle(m *telegram.NewMessage) error {
	code := m.Args()
	code, imports := resolveImports(code)

	if code == "" {
		return nil
	}

	defer os.Remove("tmp/eval.go")
	defer os.Remove("tmp/eval_out.txt")
	defer os.Remove("tmp")

	resp, isfile := perfomEval(code, m, imports)
	if isfile {
		if _, err := m.ReplyMedia(resp, telegram.MediaOptions{Caption: "Output"}); err != nil {
			m.Reply("Error: " + err.Error())
		}
		return nil
	}
	resp = strings.TrimSpace(resp)

	if resp != "" {
		if _, err := m.Reply(resp); err != nil {
			m.Reply(err)
		}
	}
	return nil
}

func perfomEval(code string, m *telegram.NewMessage, imports []string) (string, bool) {
	msg_b, _ := json.Marshal(m.Message)
	snd_b, _ := json.Marshal(m.Sender)
	cnt_b, _ := json.Marshal(m.Chat)
	chn_b, _ := json.Marshal(m.Channel)
	cache_b, _ := m.Client.Cache.ExportJSON()
	var importStatement string = ""
	if len(imports) > 0 {
		importStatement = "import (\n"
		for _, v := range imports {
			importStatement += `"` + v + `"` + "\n"
		}
		importStatement += ")\n"
	}

	code_file := fmt.Sprintf(boiler_code_for_eval, importStatement, m.ID, msg_b, snd_b, cnt_b, chn_b, cache_b, code, m.Client.ExportSession())
	tmp_dir := "tmp"
	_, err := os.ReadDir(tmp_dir)
	if err != nil {
		err = os.Mkdir(tmp_dir, 0o755)
		if err != nil {
			fmt.Println(err)
		}
	}

	// defer os.Remove(tmp_dir)

	os.WriteFile(tmp_dir+"/eval.go", []byte(code_file), 0o644)
	cmd := exec.Command("go", "run", "tmp/eval.go")
	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err = cmd.Run()
	if stdOut.String() == "" && stdErr.String() == "" {
		if err != nil {
			return fmt.Sprintf("<b>#EVALERR:</b> <code>%s</code>", err.Error()), false
		}
		return "<b>#EVALOut:</b> <code>No Output</code>", false
	}

	if stdOut.String() != "" {
		if len(stdOut.String()) > 4095 {
			os.WriteFile("tmp/eval_out.txt", stdOut.Bytes(), 0o644)
			return "tmp/eval_out.txt", true
		}

		strDou := strings.Split(stdOut.String(), "output-start")

		return fmt.Sprintf("<b>#EVALOut:</b> <code>%s</code>", strings.TrimSpace(strDou[1])), false
	}

	if stdErr.String() != "" {
		regexErr := regexp.MustCompile(`eval.go:\d+:\d+:`)
		errMsg := regexErr.Split(stdErr.String(), -1)
		if len(errMsg) > 1 {
			if len(errMsg[1]) > 4095 {
				os.WriteFile("tmp/eval_out.txt", []byte(errMsg[1]), 0o644)
				return "tmp/eval_out.txt", true
			}
			return fmt.Sprintf("<b>#EVALERR:</b> <code>%s</code>", strings.TrimSpace(errMsg[1])), false
		}
		return fmt.Sprintf("<b>#EVALERR:</b> <code>%s</code>", stdErr.String()), false
	}

	return "<b>#EVALOut:</b> <code>No Output</code>", false
}
const boilerCodeForEval = `
package main

%s

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var output string

func evalCode(bot *gotgbot.Bot, ctx *ext.Context) {
	defer func() {
		if r := recover(); r != nil {
			output = fmt.Sprintf("<b>#EVALERR:</b> <code>%v</code>", r)
		}
	}()
	
	var res string
	func() {
		defer func() {
			if r := recover(); r != nil {
				res = fmt.Sprintf("<b>#EVALERR:</b> <code>%v</code>", r)
			}
		}()
		%s
	}()

	if res == "" {
		output = "<b>#EVALOut:</b> <code>Executed Successfully</code>"
	} else {
		output = res
	}
}
`

func resolveImports(code string) (string, []string) {
	var imports []string

	code = strings.ReplaceAll(code, "package main", "")

	importsRegex := regexp.MustCompile(`import\s*([\s\S]*?)|import\s*"([\s\S]*?)"`)
	importsMatches := importsRegex.FindAllStringSubmatch(code, -1)
	for _, v := range importsMatches {
		if v[1] != "" {
			imports = append(imports, v[1])
		} else {
			imports = append(imports, v[2])
		}
	}
	code = importsRegex.ReplaceAllString(code, "")

	code = regexp.MustCompile(`func\s+main\s*\s*\{[\s\S]*?\}`).ReplaceAllString(code, "")

	return code, imports
}

func EvalHandle(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveMessage == nil || ctx.EffectiveMessage.Text == "" {
		return nil
	}

	code := strings.TrimPrefix(ctx.EffectiveMessage.Text, "/eval ")
	code, imports := resolveImports(code)

	if code == "" {
		return nil
	}

	resp, isFile := performEval(code, b, ctx, imports)
	if isFile {
		_, err := ctx.EffectiveMessage.Reply(b, "Output saved as file.", &gotgbot.SendMessageOpts{})
		return err
	}

	resp = strings.TrimSpace(resp)
	if resp != "" {
		_, err := ctx.EffectiveMessage.Reply(b, resp, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
		return err
	}
	return nil
}

func performEval(code string, b *gotgbot.Bot, ctx *ext.Context, imports []string) (string, bool) {
	msgB, _ := json.Marshal(ctx.EffectiveMessage)
	usrB, _ := json.Marshal(ctx.EffectiveUser)
	chatB, _ := json.Marshal(ctx.EffectiveChat)

	importStatement := ""
	if len(imports) > 0 {
		importStatement = "import (\n"
		for _, v := range imports {
			importStatement += `"` + v + `"` + "\n"
		}
		importStatement += ")\n"
	}

	codeFile := fmt.Sprintf(boilerCodeForEval, importStatement, code)

	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	_, err := i.Eval(codeFile)
	if err != nil {
		return fmt.Sprintf("<b>#EVALERR:</b> <code>%s</code>", err.Error()), false
	}

	v, err := i.Eval("output")
	if err != nil {
		return fmt.Sprintf("<b>#EVALERR:</b> <code>%s</code>", err.Error()), false
	}

	return v.String(), false
}