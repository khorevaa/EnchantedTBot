package SessionContext

import (
	"github.com/khorevaa/EnchantedTBot/context/internal/baseContext"
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
	"strconv"
)

var _ types.SessionContextInterface = (*SessionContext)(nil)

type Option func(ctx *SessionContext) Option

type SessionContext struct {
	baseContext.Context

	session   types.SessionStorage
	message   *tbot.Message
	messageID int
}

func (ctx *SessionContext) Options(opts ...Option) (previous Option) {

	for _, opt := range opts {
		previous = opt(ctx)
	}
	return previous
}

func WithSession(session types.SessionStorage) Option {

	return func(ctx *SessionContext) Option {

		p := ctx.session
		ctx.session = session
		return WithSession(p)
	}
}

func WithMessage(message *tbot.Message) Option {

	return func(ctx *SessionContext) Option {

		p := ctx.Message()

		chatID, _ := strconv.ParseInt(message.Chat.ID, 10, 64)
		userID := int64(message.From.ID)
		ctx.Options(WithChatID(chatID), WithUserID(userID))

		ctx.messageID = message.MessageID

		return WithMessage(p)
	}
}

func WithChatID(chatID int64) Option {

	return func(ctx *SessionContext) Option {
		p := ctx.ChatID()
		ctx.Context.Options(baseContext.WithChatID(chatID))
		return WithChatID(p)
	}

}

func WithUserID(userID int64) Option {

	return func(ctx *SessionContext) Option {
		p := ctx.UserID()
		ctx.Context.Options(baseContext.WithUserID(userID))
		return WithUserID(p)
	}

}

func WithLocal(local string) Option {

	return func(ctx *SessionContext) Option {
		p := ctx.GetLocal()
		ctx.Context.Options(baseContext.WithLocal(local))
		return WithLocal(p)
	}

}

func WithI18n(i18n types.TranslateInterface) Option {

	return func(ctx *SessionContext) Option {
		p := ctx.I18n()
		ctx.Context.Options(baseContext.WithI18n(i18n))
		return WithI18n(p)
	}

}

func WithMenu(fn types.TypeMainMenuFunc) Option {

	return func(ctx *SessionContext) Option {
		p := ctx.GetMenuMarkup()
		ctx.Context.Options(baseContext.WithMenu(fn))
		return WithMenu(p)
	}
}

func (ctx *SessionContext) Message() *tbot.Message {

	return ctx.message

}

func (ctx *SessionContext) Storage() types.SessionStorage {

	return ctx.session

}

func (ctx *SessionContext) Get() (string, map[string]interface{}) {

	return ctx.session.Get(ctx.ChatID())

}

func (ctx *SessionContext) Set(state string, data ...map[string]interface{}) {

	ctx.session.Set(ctx.ChatID(), state)

	for _, m := range data {
		ctx.session.SetData(ctx.ChatID(), m)
	}

}

func (ctx *SessionContext) Reset() {

	ctx.session.Reset(ctx.ChatID())
}
