package menu

import (
	"fmt"
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
	"sort"
	"strings"
)

const DEFAULT_MARKUP_KEY = "default"

type Row struct {
	buttons []Button
}

type Button struct {
	Text   string
	Prefix string
	Fn     types.MenuHandlerFunc
}

type Menu struct {
	rows               []Row
	registeredHandlers []string
	handlers           map[string]types.MenuHandlerFunc

	ResizeKeyboard  bool
	OneTimeKeyboard bool
	Selective       bool

	i18n types.TranslateInterface

	compiledMarkup map[string]tbot.ReplyKeyboardMarkup
}

func (m *Menu) Find(menu string) int {

	return sort.SearchStrings(m.registeredHandlers, menu)

}

func (m *Menu) Add(menu string, fn types.MenuHandlerFunc) {

	m.registeredHandlers = append(m.registeredHandlers, menu)
	m.handlers[menu] = fn

}

func (m *Menu) Contain(menu string) bool {

	return strings.EqualFold(m.registeredHandlers[m.Find(menu)], menu)

}

func (m *Menu) Run(menu string, ctx types.MessageContextInterface) {

	fn, ok := m.handlers[menu]

	if ok {
		fn(ctx)
	}

}

func (m *Menu) Keyboard(locale string) tbot.ReplyKeyboardMarkup {

	return m.getOrCompileMarkup(locale)

}

func (m *Menu) GetHandler(text string) types.MenuHandlerFunc {

	str := strings.Split(text, " ")
	prefix := str[0]

	fn, ok := m.handlers[prefix]
	if !ok {
		fn, _ = m.handlers[text]
	}

	return fn

}

func (m *Menu) IsMenuText(text string) bool {

	str := strings.Split(text, " ")
	prefix := str[0]

	_, ok := m.handlers[prefix]
	if !ok {
		_, ok = m.handlers[text]
	}

	return ok

}

func (m *Menu) getOrCompileMarkup(locale string) tbot.ReplyKeyboardMarkup {

	if m.i18n == nil {
		return m.compileMarkup(DEFAULT_MARKUP_KEY)
	}

	markup, ok := m.compiledMarkup[locale]

	if !ok {
		markup = m.compileMarkup(locale)
	}

	return markup

}

func (m *Menu) compileMarkup(local string) tbot.ReplyKeyboardMarkup {

	if markup, ok := m.compiledMarkup[local]; ok {
		return markup
	}

	useI18n := m.i18n != nil && !strings.EqualFold(local, DEFAULT_MARKUP_KEY)

	if useI18n {

		m.i18n.SetLocale(local)

	} else if markup, ok := m.compiledMarkup[DEFAULT_MARKUP_KEY]; ok {
		return markup
	}

	addButton := func(btn Button) string {

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

	for _, row := range m.rows {

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
