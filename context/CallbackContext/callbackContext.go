package CallbackContext

import (
	"github.com/khorevaa/EnchantedTBot/context/internal/baseContext"
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
	"strconv"
)

var _ types.CallbackContextInterface = (*CallbackContext)(nil)

type Option func(ctx *CallbackContext) Option

type CallbackContext struct {
	baseContext.Context
	messageID      int
	menuMarkupFunc types.TypeMainMenuFunc
	callback       *tbot.CallbackQuery
	callbackData   types.CallbackDataInterface
}

func (ctx *CallbackContext) MessageID() int {

	return ctx.messageID

}

func (ctx *CallbackContext) Options(opts ...Option) (previous Option) {

	for _, opt := range opts {
		previous = opt(ctx)
	}
	return previous
}

func WithCallbackQuery(callback *tbot.CallbackQuery) Option {

	return func(ctx *CallbackContext) Option {

		p := ctx.callback

		chatID, _ := strconv.ParseInt(callback.Message.Chat.ID, 10, 64)
		userID := getUserID(callback)
		ctx.Options(WithChatID(chatID), WithUserID(userID))
		ctx.messageID = callback.Message.MessageID
		ctx.callbackData = CallbackData(callback.Data)

		return WithCallbackQuery(p)
	}

}

func getUserID(callback *tbot.CallbackQuery) int64 {

	switch {
	case callback.Message.From != nil && !callback.Message.From.IsBot:
		return int64(callback.Message.From.ID)
	case callback.From != nil && !callback.From.IsBot:
		return int64(callback.From.ID)
	case len(callback.Message.Chat.ID) > 0:
		userID, _ := strconv.ParseInt(callback.Message.Chat.ID, 10, 64)
		return userID
	default:
		return -1
	}

}

func WithChatID(chatID int64) Option {

	return func(ctx *CallbackContext) Option {
		p := ctx.ChatID()
		ctx.Context.Options(baseContext.WithChatID(chatID))
		return WithChatID(p)
	}

}

func WithUserID(userID int64) Option {

	return func(ctx *CallbackContext) Option {
		p := ctx.UserID()
		ctx.Context.Options(baseContext.WithUserID(userID))
		return WithUserID(p)
	}

}

func WithLocal(local string) Option {

	return func(ctx *CallbackContext) Option {
		p := ctx.GetLocal()
		ctx.Context.Options(baseContext.WithLocal(local))
		return WithLocal(p)
	}

}

func WithI18n(i18n types.TranslateInterface) Option {

	return func(ctx *CallbackContext) Option {
		p := ctx.I18n()
		ctx.Context.Options(baseContext.WithI18n(i18n))
		return WithI18n(p)
	}

}

func WithMenu(fn types.TypeMainMenuFunc) Option {

	return func(ctx *CallbackContext) Option {
		p := ctx.GetMenuMarkup()
		ctx.Context.Options(baseContext.WithMenu(fn))
		return WithMenu(p)
	}
}

func (ctx *CallbackContext) Callback() *tbot.CallbackQuery {

	return ctx.callback

}

func (ctx *CallbackContext) CallbackData() types.CallbackDataInterface {

	return ctx.callbackData

}

func (ctx *CallbackContext) Next(n types.CallbackActionInterface, d ...string) string {

	return ctx.callbackData.Next(n, d...)

}

func (ctx *CallbackContext) Back(d ...string) string {

	return ctx.callbackData.Back(d...)

}
