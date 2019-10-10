package baseContext

import (
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
)

var _ types.ContextInterface = (*Context)(nil)

type Option func(ctx *Context) Option

type Context struct {
	userID          int64
	chatID          int64
	menuMarkupFunc  types.TypeMainMenuFunc
	i18n            types.TranslateInterface
	localTranslator types.TranslateServiceInterface
}

func (ctx *Context) Options(opts ...Option) (previous Option) {

	for _, opt := range opts {
		previous = opt(ctx)
	}
	return previous
}
func WithChatID(chatID int64) Option {

	return func(ctx *Context) Option {
		p := ctx.chatID
		ctx.chatID = chatID
		return WithChatID(p)
	}

}

func WithUserID(userID int64) Option {

	return func(ctx *Context) Option {
		p := ctx.userID
		ctx.userID = userID
		return WithUserID(p)
	}

}

func WithLocal(local string) Option {

	return func(ctx *Context) Option {
		p := ctx.GetLocal()
		ctx.SetLocal(local)
		return WithLocal(p)
	}

}

func WithI18n(i18n types.TranslateInterface) Option {

	return func(ctx *Context) Option {
		p := ctx.i18n
		ctx.i18n = i18n
		return WithI18n(p)
	}

}

func WithMenu(fn types.TypeMainMenuFunc) Option {

	return func(ctx *Context) Option {
		p := ctx.menuMarkupFunc
		ctx.menuMarkupFunc = fn
		return WithMenu(p)
	}

}

func (ctx *Context) SetMenuMarkup(f types.TypeMainMenuFunc) {
	ctx.menuMarkupFunc = f
}

func (ctx *Context) GetMenuMarkup() types.TypeMainMenuFunc {
	return ctx.menuMarkupFunc
}

func (ctx *Context) SetI18n(f types.TranslateInterface) {
	ctx.i18n = f
}

func (ctx *Context) I18n() types.TranslateInterface {
	return ctx.i18n
}

func (ctx *Context) UserID() int64 {

	return ctx.userID

}

func (ctx *Context) ChatID() int64 {

	return ctx.chatID

}

func (ctx *Context) MenuMarkup() tbot.ReplyKeyboardMarkup {

	fn := ctx.menuMarkupFunc

	if fn == nil {
		return tbot.ReplyKeyboardMarkup{}
	}

	return fn(ctx)

}

func (ctx *Context) GetLocal() string {

	return ctx.i18n.GetLocal()

}

func (ctx *Context) SetLocal(local string) types.TranslateServiceInterface {

	ctx.localTranslator = ctx.i18n.SetLocal(local)

	return ctx.localTranslator

}

func (ctx *Context) F(format string, args ...interface{}) string {

	return ctx.i18n.F(format, args)
}

func (ctx *Context) T(key string) string {

	return ctx.i18n.T(key)
}

func (ctx *Context) NStr(m interface{}, args ...interface{}) string {

	return ctx.i18n.NStr(m, args...)
}
