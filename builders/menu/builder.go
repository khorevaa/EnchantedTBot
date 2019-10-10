package menu

import (
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
)

type i18n interface {
	SetLocale(string)
	T(string) string
}

type ContextMenuInterface interface {
}

type Option func(*MenuBuilder) Option

type MenuBuilder struct {
	menu            []*buttonsRow
	currentRow      *buttonsRow
	ResizeKeyboard  bool `json:"resize_keyboard"`
	OneTimeKeyboard bool `json:"one_time_keyboard"`
	Selective       bool `json:"selective"`

	i18n i18n

	mainMenu *Menu
}

func (f *MenuBuilder) Option(opts ...Option) {
	for _, opt := range opts {
		opt(f)
	}
}

type buttonsRow struct {
	buttons []menuButton
}

type menuButton struct {
	Text   string
	Prefix string
	Fn     types.MenuHandlerFunc
}

func NewMenuBuilder(opts ...Option) *MenuBuilder {

	b := &MenuBuilder{}

	b.Option(opts...)

	b.initFirstRow()

	return b
}

func (b *MenuBuilder) initFirstRow() {

	row := &buttonsRow{}
	b.menu = append(b.menu, row)
	b.currentRow = row

}

func OneTimeKeyboard(oneTime bool) Option {

	return func(b *MenuBuilder) Option {

		p := b.OneTimeKeyboard
		b.OneTimeKeyboard = oneTime

		return OneTimeKeyboard(p)
	}
}

func WithTranslate(i18n i18n) func(b *MenuBuilder) {

	return func(b *MenuBuilder) {
		b.i18n = i18n
	}

}

func WithMainMenu(mm *Menu) Option {

	return func(b *MenuBuilder) Option {

		p := b.mainMenu
		b.mainMenu = mm

		return WithMainMenu(p)
	}

}

func ResizeKeyboard(resize bool) Option {

	return func(b *MenuBuilder) Option {

		p := b.ResizeKeyboard
		b.ResizeKeyboard = resize

		return ResizeKeyboard(p)
	}

}
func SelectiveKeyboard(selective bool) Option {

	return func(b *MenuBuilder) Option {

		p := b.Selective
		b.Selective = selective

		return SelectiveKeyboard(p)
	}

}

func newButton(text string, fn types.MenuHandlerFunc, opts ...func(btn menuButton)) menuButton {

	btn := menuButton{
		Text: text,
		Fn:   fn,
	}

	for _, f := range opts {
		f(btn)
	}

	return btn
}

func (b *MenuBuilder) Button(text string, fn types.MenuHandlerFunc, opts ...func(btn menuButton)) *MenuBuilder {

	b.currentRow.buttons = append(b.currentRow.buttons, newButton(text, fn, opts...))

	return b
}

func (b *MenuBuilder) NewRow() *MenuBuilder {

	row := &buttonsRow{}
	b.menu = append(b.menu, row)
	b.currentRow = row

	return b
}

func (b *MenuBuilder) Build() *Menu {

	mainMenu := &Menu{
		handlers:       make(map[string]types.MenuHandlerFunc),
		compiledMarkup: make(map[string]tbot.ReplyKeyboardMarkup),
	}

	mainMenu.rows = b.menu
	mainMenu.Selective = b.Selective
	mainMenu.OneTimeKeyboard = b.OneTimeKeyboard
	mainMenu.ResizeKeyboard = b.ResizeKeyboard

	mainMenu.compileMarkup(DEFAULT_MARKUP_KEY)

	return mainMenu
}

func WithEmoji(emoji string) func(btn menuButton) {

	return WithPrefix(emoji)

}

func WithPrefix(prefix string) func(btn menuButton) {

	return func(btn menuButton) {
		btn.Prefix = prefix
	}

}

func (b *MenuBuilder) ButtonWithEmoji(text string, fn types.MenuHandlerFunc, emoji string) *MenuBuilder {

	b.currentRow.buttons = append(b.currentRow.buttons, newButton(text, fn, WithEmoji(emoji)))

	return b
}
