package types

import (
	"github.com/yanzay/tbot/v2"
)

type StateHandlerFunc func(ctx SessionContextInterface)
type CallbackHandlerFunc func(ctx CallbackContextInterface)
type MenuHandlerFunc func(ctx MessageContextInterface)

type CommandHandlerFunc func(ctx CommandContextInterface)

type TypeMainMenuFunc func(ctx ContextInterface) tbot.ReplyKeyboardMarkup

type SessionStorage interface {
	Set(int64, string)
	Get(int64) (string, map[string]interface{})
	SetData(int64, map[string]interface{})
	GetData(int64) map[string]interface{}
	Reset(int64)
}

type ContextInterface interface {
	UserID() int64
	ChatID() int64

	MenuMarkup() tbot.ReplyKeyboardMarkup

	SetMenuMarkup(TypeMainMenuFunc)
	SetI18n(TranslateInterface)
}

//type ContextOptionsInterface interface {
//	WithLocal(ContextOptionsInterface) interface{}
//}

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

type TranslateServiceInterface interface {
	GetLocal() string
	F(format string, args ...interface{}) string
	T(key string) string
	NStr(m interface{}, args ...interface{}) string
}

type TranslateInterface interface {
	SetLocal(local string) TranslateServiceInterface
	TranslateServiceInterface
}

type MenuMarkupFuncInterface interface {
	MenuMarkupFunc() TypeMainMenuFunc
}

type RegisterHandlersInterface interface {
	CallbackHandler(callback CallbackActionInterface, h CallbackHandlerFunc)
	StateHandler(state string, h StateHandlerFunc)
	MenuHandler(menu string, h MenuHandlerFunc)
}

type CallbackDataInterface interface {
	Next(next CallbackActionInterface, data ...string) string
	Back(data ...string) string
	Action() CallbackActionInterface
	Data() CallbackActionDataInterface
	WithData(data ...string) CallbackDataInterface
}

type CallbackActionInterface interface {
	Value() int8
	String() string
}

type CallbackActionDataInterface interface {
	Value() string
	String() string
	Map() map[string]string
	FromSlice(args ...string) CallbackActionDataInterface
	FromMap(in map[string]string) CallbackActionDataInterface
}

type HandlerContextInterface interface {
	UserID() int64
	ChatID() int64

	MenuMarkup() tbot.ReplyKeyboardMarkup

	TranslateInterface
}
