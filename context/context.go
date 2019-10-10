package context

import (
	"github.com/khorevaa/EnchantedTBot/context/CallbackContext"
	"github.com/khorevaa/EnchantedTBot/context/CommandContext"
	"github.com/khorevaa/EnchantedTBot/context/MessageContext"
	"github.com/khorevaa/EnchantedTBot/context/SessionContext"
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
)

func NewCallbackContext(callback *tbot.CallbackQuery, opts ...CallbackContext.Option) types.CallbackContextInterface {

	ctx := &CallbackContext.CallbackContext{}
	ctx.Options(CallbackContext.WithCallbackQuery(callback))
	ctx.Options(opts...)

	return ctx

}

func NewMessageContext(message *tbot.Message, opts ...MessageContext.Option) types.MessageContextInterface {

	ctx := &MessageContext.MessageContext{}
	ctx.Options(MessageContext.WithMessage(message))
	ctx.Options(opts...)

	return ctx

}

func NewCommandContext(message *tbot.Message, opts ...CommandContext.Option) types.CommandContextInterface {

	ctx := &CommandContext.CommandContext{}
	ctx.Options(CommandContext.WithMessage(message))
	ctx.Options(opts...)

	return ctx

}

func NewStateContext(message *tbot.Message, session types.SessionStorage, opts ...SessionContext.Option) types.SessionContextInterface {

	ctx := &SessionContext.SessionContext{}
	ctx.Options(SessionContext.WithMessage(message), SessionContext.WithSession(session))
	ctx.Options(opts...)

	return ctx

}
