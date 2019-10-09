package EnchantedTBot

import (
	"fmt"
	"github.com/yanzay/tbot/v2"
	"strings"
)

const DEFAULT_MARKUP_KEY = "default"

type mainMenu struct {
	menu []buttonsRow

	handlers map[string]func(ContextMenuInterface)

	ResizeKeyboard  bool `json:"resize_keyboard"`
	OneTimeKeyboard bool `json:"one_time_keyboard"`
	Selective       bool `json:"selective"`

	i18n i18n

	compiledMarkup map[string]tbot.ReplyKeyboardMarkup
}

func (m *mainMenu) Keyboard(locale string) tbot.ReplyKeyboardMarkup {

	return m.getOrCompileMarkup(locale)

}

func (m *mainMenu) GetHandler(text string) func(ContextMenuInterface) {

	str := strings.Split(text, " ")
	prefix := str[0]

	fn, ok := m.handlers[prefix]
	if !ok {
		fn, _ = m.handlers[text]
	}

	return fn

}

func (m *mainMenu) IsMenuText(text string) bool {

	str := strings.Split(text, " ")
	prefix := str[0]

	_, ok := m.handlers[prefix]
	if !ok {
		_, ok = m.handlers[text]
	}

	return ok

}

func (m *mainMenu) getOrCompileMarkup(locale string) tbot.ReplyKeyboardMarkup {

	if m.i18n == nil {
		return m.compileMarkup(DEFAULT_MARKUP_KEY)
	}

	markup, ok := m.compiledMarkup[locale]

	if !ok {
		markup = m.compileMarkup(locale)
	}

	return markup

}

func (m *mainMenu) compileMarkup(local string) tbot.ReplyKeyboardMarkup {

	if markup, ok := m.compiledMarkup[local]; ok {
		return markup
	}

	useI18n := m.i18n != nil && !strings.EqualFold(local, DEFAULT_MARKUP_KEY)

	if useI18n {

		m.i18n.SetLocale(local)

	} else if markup, ok := m.compiledMarkup[DEFAULT_MARKUP_KEY]; ok {
		return markup
	}

	addButton := func(btn menuButton) string {

		btnText := btn.Text

		if useI18n {

			btnText = m.i18n.T(btn.Text)

		}

		handlerKey := btnText

		if len(btn.Prefix) > 0 {
			handlerKey = btn.Prefix
			btnText = fmt.Sprintf("%s %s", btn.Prefix, btnText)
		}

		m.handlers[handlerKey] = btn.Fn

		return btnText

	}

	markup := tbot.ReplyKeyboardMarkup{
		ResizeKeyboard:  m.ResizeKeyboard,
		OneTimeKeyboard: m.OneTimeKeyboard,
		Selective:       m.Selective,
	}

	for _, row := range m.menu {

		var btnRow []tbot.KeyboardButton

		for _, btn := range row.buttons {

			btnText := addButton(btn)

			btnRow = append(btnRow, tbot.KeyboardButton{
				Text: btnText,
			})
		}

		markup.Keyboard = append(markup.Keyboard, btnRow)

	}

	m.compiledMarkup[local] = markup

	return markup
}

type i18n interface {
	SetLocale(string)
	T(string) string
}

type ContextMenuInterface interface {
}

type MenuBuilder struct {
	menu            []buttonsRow
	currentRow      *buttonsRow
	ResizeKeyboard  bool `json:"resize_keyboard"`
	OneTimeKeyboard bool `json:"one_time_keyboard"`
	Selective       bool `json:"selective"`

	i18n i18n

	mainMenu *mainMenu
}

type buttonsRow struct {
	buttons []menuButton
}

type menuButton struct {
	Text   string
	Prefix string
	Fn     func(ContextMenuInterface)
}

func NewMenuBuilder(opts ...func(b *MenuBuilder)) *MenuBuilder {

	b := &MenuBuilder{}

	for _, f := range opts {
		f(b)
	}

	return b
}

func OneTimeKeyboard(oneTime bool) func(b *MenuBuilder) {

	return func(b *MenuBuilder) {
		b.OneTimeKeyboard = oneTime
	}

}

func WithTranslate(i18n i18n) func(b *MenuBuilder) {

	return func(b *MenuBuilder) {
		b.i18n = i18n
	}

}

func ResizeKeyboard(resize bool) func(b *MenuBuilder) {

	return func(b *MenuBuilder) {
		b.ResizeKeyboard = resize
	}

}
func SelectiveKeyboard(selective bool) func(b *MenuBuilder) {

	return func(b *MenuBuilder) {
		b.Selective = selective
	}

}

func NewButton(text string, fn func(ContextMenuInterface), opts ...func(btn menuButton)) menuButton {

	btn := menuButton{
		Text: text,
		Fn:   fn,
	}

	for _, f := range opts {
		f(btn)
	}

	return btn
}

func (b *MenuBuilder) Button(text string, fn func(ContextMenuInterface), opts ...func(btn menuButton)) *MenuBuilder {

	b.currentRow.buttons = append(b.currentRow.buttons, NewButton(text, fn, opts...))

	return b
}

func (b *MenuBuilder) NewRow() *MenuBuilder {

	row := buttonsRow{}
	b.menu = append(b.menu, row)
	b.currentRow = &row

	return b
}

func (b *MenuBuilder) Build() *mainMenu {

	b.mainMenu.menu = b.menu
	b.mainMenu.Selective = b.Selective
	b.mainMenu.OneTimeKeyboard = b.OneTimeKeyboard
	b.mainMenu.ResizeKeyboard = b.ResizeKeyboard

	return b
}

func WithEmoji(emoji string) func(btn menuButton) {

	return WithPrefix(emoji)

}

func WithPrefix(prefix string) func(btn menuButton) {

	return func(btn menuButton) {
		btn.Prefix = prefix
	}

}

func (b *MenuBuilder) ButtonWithEmoji(text string, emoji string) *MenuBuilder {

	b.currentRow.buttons = append(b.currentRow.buttons, NewButton(text, WithEmoji(emoji)))

	return b
}
