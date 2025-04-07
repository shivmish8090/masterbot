package helper

import (
	"github.com/PaulSonOfLars/gotgbot/v2/telegram"
)

// Buttons is a smart inline button builder
type Buttons struct {
	RowWidth int
	rows     [][]telegram.InlineKeyboardButton
	tmpRow   []telegram.InlineKeyboardButton
}

// Inline creates a callback button
func (b *Buttons) Inline(text, data string) telegram.InlineKeyboardButton {
	return telegram.InlineKeyboardButton{
		Text:         text,
		CallbackData: data,
	}
}

// Url creates a url button
func (b *Buttons) Url(text, url string) telegram.InlineKeyboardButton {
	return telegram.InlineKeyboardButton{
		Text: text,
		Url:  url,
	}
}

// Add adds buttons automatically using RowWidth rule
func (b *Buttons) Add(buttons ...telegram.InlineKeyboardButton) {
	for _, btn := range buttons {
		b.tmpRow = append(b.tmpRow, btn)
		if len(b.tmpRow) == b.RowWidth {
			b.rows = append(b.rows, b.tmpRow)
			b.tmpRow = nil
		}
	}
}

// Row forces a new row with given buttons (ignores RowWidth)
func (b *Buttons) Row(buttons ...telegram.InlineKeyboardButton) {
	if len(b.tmpRow) > 0 {
		b.rows = append(b.rows, b.tmpRow)
		b.tmpRow = nil
	}
	b.rows = append(b.rows, buttons)
}

// Build builds InlineKeyboardMarkup final structure
func (b *Buttons) Build() telegram.InlineKeyboardMarkup {
	if len(b.tmpRow) > 0 {
		b.rows = append(b.rows, b.tmpRow)
		b.tmpRow = nil
	}
	return telegram.InlineKeyboardMarkup{InlineKeyboard: b.rows}
}
