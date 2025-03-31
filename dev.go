package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func LsHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	text := ctx.EffectiveMessage.GetText()
	fields := strings.Fields(text)
	dir := ""
	if len(fields) < 1 {
		dir = "."
	}
	dir = strings.TrimSpace(strings.Replace(text, fields[0], "", 1))

	cmd := exec.Command("ls", dir)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
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

	if err != nil {
		ctx.EffectiveMessage.Reply(b, "<code>Error:</code> <b>"+err.Error()+"</b>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})
		return nil
	}

	files := strings.Split(strings.TrimSpace(out.String()), "\n")
	var sizeTotal int64

	var resp string
	for _, file := range files {
		fileType := "file"
		if strings.Contains(file, ".") {
			fp := strings.Split(file, ".")
			fileType = fp[len(fp)-1]
		}
		switch fileType {
		case "mp4", "mkv", "webm", "avi", "flv", "mov", "wmv", "3gp":
			fileType = "video"
		case "mp3", "wav", "flac", "ogg", "m4a", "wma":
			fileType = "audio"
		case "jpg", "jpeg", "png", "gif", "webp", "bmp", "tiff":
			fileType = "image"
		case "go":
			fileType = "go"
		case "py":
			fileType = "python"
		case "txt":
			fileType = "txt"
		default:
			fileType = "file"
		}
		size := calcFileOrDirSize(filepath.Join(dir, file))
		sizeTotal += size
		resp += fileTypeEmoji[fileType] + " " + file + " " + "(" + sizeToHuman(size) + ")" + "\n"
	}

	resp += "\nTotal: " + sizeToHuman(sizeTotal)

	ctx.EffectiveMessage.Reply(b, "<pre lang='bash'>"+resp+"</pre>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})

	return nil
}

func sizeToHuman(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
}

func calcFileOrDirSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}

	if !fi.IsDir() {
		return fi.Size()
	}

	var size int64
	walker := func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fi, err := info.Info()
			if err != nil {
				return err
			}
			size += fi.Size()
		}
		return nil
	}

	err = filepath.WalkDir(path, walker)
	if err != nil {
		return 0
	}

	return size
}
