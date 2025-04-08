package utils

import "github.com/PaulSonOfLars/gotgbot/v2"

// Buttons is a smart inline button builder
type Buttons struct {
	RowWidth int
	rows     [][]gotgbot.InlineKeyboardButton
	tmpRow   []gotgbot.InlineKeyboardButton
}

// Inline creates a callback button
func (b *Buttons) Inline(text, data string) gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text:         text,
		CallbackData: data,
	}
}

// Url creates a url button
func (b *Buttons) Url(text, url string) gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text: text,
		Url:  url,
	}
}

// Add adds buttons automatically using RowWidth rule
func (b *Buttons) Add(buttons ...gotgbot.InlineKeyboardButton) {
	for _, btn := range buttons {
		b.tmpRow = append(b.tmpRow, btn)
		if len(b.tmpRow) == b.RowWidth {
			b.rows = append(b.rows, b.tmpRow)
			b.tmpRow = nil
		}
	}
}

// Row forces a new row with given buttons (ignores RowWidth)
func (b *Buttons) Row(buttons ...gotgbot.InlineKeyboardButton) {
	if len(b.tmpRow) > 0 {
		b.rows = append(b.rows, b.tmpRow)
		b.tmpRow = nil
	}
	b.rows = append(b.rows, buttons)
}

// Build builds InlineKeyboardMarkup final structure
func (b *Buttons) Build() gotgbot.InlineKeyboardMarkup {
	if len(b.tmpRow) > 0 {
		b.rows = append(b.rows, b.tmpRow)
		b.tmpRow = nil
	}
	return gotgbot.InlineKeyboardMarkup{InlineKeyboard: b.rows}
}
