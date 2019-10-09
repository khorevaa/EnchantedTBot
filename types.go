package HBot

import (
	"github.com/yanzay/tbot/v2"
)

type StateHandlerFunc func(ctx SessionContextInterface)
type CallbackHandlerFunc func(ctx CallbackContextInterface)
type MenuHandlerFunc func(ctx MessageContextInterface)

type CommandHandlerFunc func(ctx CommandContextInterface)

type TypeMainMenuFunc func(ctx ContextInterface) tbot.ReplyKeyboardMarkup

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
