package HBot

import (
	"github.com/yanzay/tbot/v2"
	"strconv"
)

var _ ContextInterface = (*baseContext)(nil)
var _ MessageContextInterface = (*MessageContext)(nil)
var _ CallbackContextInterface = (*CallbackContext)(nil)
var _ CommandContextInterface = (*CommandContext)(nil)
var _ SessionContextInterface = (*SessionContext)(nil)

type ContextInterface interface {
	UserID() int64
	ChatID() int64

	MenuMarkup() tbot.ReplyKeyboardMarkup

	SetMenuMarkup(TypeMainMenuFunc)
	SetI18n(TranslateInterface)
}

type MessageContextInterface interface {
	ContextInterface

	MessageID() int
	Text() string
	Message() *tbot.Message
}

type CommandContextInterface interface {
	ContextInterface

	Command() string
	CommandArguments() []string
	Message() *tbot.Message
}

type CallbackContextInterface interface {
	ContextInterface

	MessageID() int
	CallbackData() CallbackDataInterface
	Callback() *tbot.CallbackQuery
	Next(CallbackActionInterface, ...string) string
	Back(...string) string
}

type SessionContextInterface interface {
	Message() *tbot.Message
	Storage() SessionStorage
	Set(string, ...map[string]interface{})
	Get() (string, map[string]interface{})
	Reset()
}

type baseContext struct {
	userID          int64
	chatID          int64
	menuMarkupFunc  TypeMainMenuFunc
	i18n            TranslateInterface
	localTranslator TranslateServiceInterface
}

type MessageContext struct {
	baseContext
	message   *tbot.Message
	messageID int
}

type CommandContext struct {
	baseContext

	command string
	args    []string
	message *tbot.Message
}

type CallbackContext struct {
	baseContext
	messageID      int
	menuMarkupFunc TypeMainMenuFunc
	callback       *tbot.CallbackQuery
	callbackData   CallbackDataInterface
}

type SessionContext struct {
	baseContext

	session SessionStorage
	message *tbot.Message
}

func NewCallbackContext(callback *tbot.CallbackQuery, options ...func(ContextInterface)) CallbackContextInterface {

	ctx := &CallbackContext{
		callback: callback,
	}

	ctx.chatID, _ = strconv.ParseInt(callback.Message.Chat.ID, 10, 64)
	ctx.messageID = callback.Message.MessageID
	ctx.callbackData = CallbackData(callback.Data)

	switch {

	case callback.Message.From != nil && !callback.Message.From.IsBot:
		ctx.userID = int64(callback.Message.From.ID)
	case callback.From != nil && !callback.From.IsBot:
		ctx.userID = int64(callback.From.ID)
	case len(callback.Message.Chat.ID) > 0:
		ctx.userID, _ = strconv.ParseInt(callback.Message.Chat.ID, 10, 64)
	default:
		ctx.userID = -1
	}

	for _, opt := range options {
		opt(ctx)
	}

	return ctx

}

func NewMessageContext(message *tbot.Message, options ...func(ContextInterface)) MessageContextInterface {

	ctx := &MessageContext{
		message: message,
	}

	ctx.chatID, _ = strconv.ParseInt(message.Chat.ID, 10, 64)
	ctx.messageID = message.MessageID
	ctx.userID = int64(message.From.ID)
	for _, opt := range options {
		opt(ctx)
	}

	return ctx

}

func NewCommandContext(message *tbot.Message, options ...func(ContextInterface)) CommandContextInterface {

	ctx := &CommandContext{
		message: message,
	}

	ctx.chatID, _ = strconv.ParseInt(message.Chat.ID, 10, 64)
	ctx.userID = int64(message.From.ID)

	for _, opt := range options {
		opt(ctx)
	}

	return ctx

}

func NewStateContext(message *tbot.Message, session SessionStorage, options ...func(ContextInterface)) SessionContextInterface {

	ctx := &SessionContext{
		message: message,
		session: session,
	}

	ctx.chatID, _ = strconv.ParseInt(message.Chat.ID, 10, 64)
	ctx.userID = int64(message.From.ID)

	for _, opt := range options {
		opt(ctx)
	}

	return ctx

}

func (ctx *CallbackContext) MessageID() int {

	return ctx.messageID

}

func (ctx *CallbackContext) Callback() *tbot.CallbackQuery {

	return ctx.callback

}

func (ctx *CallbackContext) CallbackData() CallbackDataInterface {

	return ctx.callbackData

}

func (ctx *CallbackContext) Next(n CallbackActionInterface, d ...string) string {

	return ctx.callbackData.Next(n, d...)

}

func (ctx *CallbackContext) Back(d ...string) string {

	return ctx.callbackData.Back(d...)

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

func (ctx *CommandContext) Message() *tbot.Message {

	return ctx.message

}

func (ctx *CommandContext) CommandArguments() []string {

	return ctx.args

}

func (ctx *CommandContext) Command() string {

	return ctx.command

}

func (ctx *SessionContext) Message() *tbot.Message {

	return ctx.message

}

func (ctx *SessionContext) Storage() SessionStorage {

	return ctx.session

}

func (ctx *SessionContext) Get() (string, map[string]interface{}) {

	return ctx.session.Get(ctx.chatID), ctx.session.GetData(ctx.chatID)

}

func (ctx *SessionContext) Set(state string, data ...map[string]interface{}) {

	ctx.session.Set(ctx.chatID, state)

	for _, m := range data {
		ctx.session.SetData(ctx.chatID, m)
	}

}

func (ctx *SessionContext) Reset() {

	ctx.session.Reset(ctx.chatID)
}

func (ctx *baseContext) SetMenuMarkup(f TypeMainMenuFunc) {
	ctx.menuMarkupFunc = f
}

func (ctx *baseContext) SetI18n(f TranslateInterface) {
	ctx.i18n = f
}

func (ctx *baseContext) UserID() int64 {

	return ctx.userID

}

func (ctx *baseContext) ChatID() int64 {

	return ctx.chatID

}

func (ctx *baseContext) MenuMarkup() tbot.ReplyKeyboardMarkup {

	fn := ctx.menuMarkupFunc

	if fn == nil {
		return tbot.ReplyKeyboardMarkup{}
	}

	return fn(ctx)

}

func (ctx *baseContext) GetLocal() string {

	return ctx.i18n.GetLocal()

}

func (ctx *baseContext) SetLocal(local string) TranslateServiceInterface {

	ctx.localTranslator = ctx.i18n.SetLocal(local)

	return ctx.localTranslator

}

func (ctx *baseContext) F(format string, args ...interface{}) string {

	return ctx.localTranslator.F(format, args)
}

func (ctx *baseContext) T(key string) string {

	return ctx.localTranslator.T(key)
}

func (ctx *baseContext) NStr(m interface{}, args ...interface{}) string {

	return ctx.localTranslator.NStr(m, args...)
}
