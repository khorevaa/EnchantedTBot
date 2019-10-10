package MessageContext

import (
	"github.com/khorevaa/EnchantedTBot/context/internal/baseContext"
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
	"strconv"
)

var _ types.MessageContextInterface = (*MessageContext)(nil)

type Option func(ctx *MessageContext) Option

type MessageContext struct {
	baseContext.Context
	message   *tbot.Message
	messageID int
}

func (ctx *MessageContext) Options(opts ...Option) (previous Option) {

	for _, opt := range opts {
		previous = opt(ctx)
	}
	return previous
}

func WithMessage(message *tbot.Message) Option {

	return func(ctx *MessageContext) Option {

		p := ctx.Message()

		chatID, _ := strconv.ParseInt(message.Chat.ID, 10, 64)
		userID := int64(message.From.ID)
		ctx.Options(WithChatID(chatID), WithUserID(userID))

		ctx.messageID = message.MessageID

		return WithMessage(p)
	}

}

func WithChatID(chatID int64) Option {

	return func(ctx *MessageContext) Option {
		p := ctx.ChatID()
		ctx.Context.Options(baseContext.WithChatID(chatID))
		return WithChatID(p)
	}

}

func WithUserID(userID int64) Option {

	return func(ctx *MessageContext) Option {
		p := ctx.UserID()
		ctx.Context.Options(baseContext.WithUserID(userID))
		return WithUserID(p)
	}

}

func WithLocal(local string) Option {

	return func(ctx *MessageContext) Option {
		p := ctx.GetLocal()
		ctx.Context.Options(baseContext.WithLocal(local))
		return WithLocal(p)
	}

}

func WithI18n(i18n types.TranslateInterface) Option {

	return func(ctx *MessageContext) Option {
		p := ctx.I18n()
		ctx.Context.Options(baseContext.WithI18n(i18n))
		return WithI18n(p)
	}

}

func WithMenu(fn types.TypeMainMenuFunc) Option {

	return func(ctx *MessageContext) Option {
		p := ctx.GetMenuMarkup()
		ctx.Context.Options(baseContext.WithMenu(fn))
		return WithMenu(p)
	}
}

func (ctx *MessageContext) MessageID() int {

	return ctx.messageID

}

func (ctx *MessageContext) Message() *tbot.Message {

	return ctx.message

}

func (ctx *MessageContext) Text() string {

	return ctx.message.Text

}
