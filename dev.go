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
