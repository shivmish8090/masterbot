package buttons

import "github.com/PaulSonOfLars/gotgbot/v2"

type Button struct {
        RowWidth int
        rows     [][]gotgbot.InlineKeyboardButton
        tmpRow   []gotgbot.InlineKeyboardButton
}

func (b *Button) Inline(text, data string) gotgbot.InlineKeyboardButton {
        return gotgbot.InlineKeyboardButton{
                Text:         text,
                CallbackData: data,
        }
}

func (b *Button) Url(text, url string) gotgbot.InlineKeyboardButton {
        return gotgbot.InlineKeyboardButton{
                Text: text,
                Url:  url,
        }
}

func (b *Button) Add(buttons ...gotgbot.InlineKeyboardButton) {
        for _, btn := range buttons {
                b.tmpRow = append(b.tmpRow, btn)
                if len(b.tmpRow) == b.RowWidth {
                        b.rows = append(b.rows, b.tmpRow)
                        b.tmpRow = nil
                }
        }
}

func (b *Button) Row(buttons ...gotgbot.InlineKeyboardButton) {
        if len(b.tmpRow) > 0 {
                b.rows = append(b.rows, b.tmpRow)
                b.tmpRow = nil
        }
        b.rows = append(b.rows, buttons)
}

func (b *Button) Build() gotgbot.InlineKeyboardMarkup {
        if len(b.tmpRow) > 0 {
                b.rows = append(b.rows, b.tmpRow)
                b.tmpRow = nil
        }
        return gotgbot.InlineKeyboardMarkup{InlineKeyboard: b.rows}
}