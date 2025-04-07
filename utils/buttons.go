package utils

type Buttons struct {
	RowWidth int
	rows     [][]telegram.InlineKeyboardButton
	tmpRow   []telegram.InlineKeyboardButton
}

func (b *Buttons) Inline(text, data string) telegram.InlineKeyboardButton {
	return telegram.InlineKeyboardButton{
		Text:         text,
		CallbackData: data,
	}
}

func (b *Buttons) Url(text, url string) telegram.InlineKeyboardButton {
	return telegram.InlineKeyboardButton{
		Text: text,
		Url:  url,
	}
}

// Add auto arrange buttons with RowWidth logic
func (b *Buttons) Add(buttons ...telegram.InlineKeyboardButton) {
	for _, btn := range buttons {
		b.tmpRow = append(b.tmpRow, btn)
		if len(b.tmpRow) == b.RowWidth {
			b.rows = append(b.rows, b.tmpRow)
			b.tmpRow = nil
		}
	}
}

// Row forcefully adds a new row (no RowWidth limit here)
func (b *Buttons) Row(buttons ...telegram.InlineKeyboardButton) {
	if len(b.tmpRow) > 0 {
		b.rows = append(b.rows, b.tmpRow)
		b.tmpRow = nil
	}
	b.rows = append(b.rows, buttons)
}

// Build final keyboard
func (b *Buttons) Build() telegram.InlineKeyboardMarkup {
	if len(b.tmpRow) > 0 {
		b.rows = append(b.rows, b.tmpRow)
		b.tmpRow = nil
	}
	return telegram.InlineKeyboardMarkup{InlineKeyboard: b.rows}
}
