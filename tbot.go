package HBot

import (

	"github.com/khorevaa/EnchantedTBot/csm"
	"github.com/yanzay/tbot/v2"
)

type BotWithHandlers struct {
	*tbot.Server

	StateHandlers    map[string]StateHandlerFunc
	MenuHandlers     map[string]MenuHandlerFunc
	CallbackHandlers map[int8]CallbackHandlerFunc
	CommandHandlers  map[string]CommandHandlerFunc

	Sessions SessionStorage

	MainMenuFunc TypeMainMenuFunc

	ErrorCallbackHandler func(callback CallbackDataInterface, ctx HandlerContextInterface)
}

func (h *BotWithHandlers) InitHandlers() {

	h.StateHandlers = make(map[string]StateHandlerFunc)
	h.MenuHandlers = make(map[string]MenuHandlerFunc)
	h.CallbackHandlers = make(map[int8]CallbackHandlerFunc)
	h.CommandHandlers = make(map[string]CommandHandlerFunc)

}

func (h *BotWithHandlers) registerCallbackHandler(callback CallbackActionInterface, handler CallbackHandlerFunc) {

	key := callback.Value()
	h.CallbackHandlers[key] = handler

}

func (h *BotWithHandlers) registerStateHandler(state string, handler StateHandlerFunc) {

	if h.GetChatStateMachine() == nil {
		h.RegisterChatStateMachine(&csm.InMenoryCsm{})
	}

	h.StateHandlers[state] = handler

}

func (h *BotWithHandlers) RegisterMenuFunc(m TypeMainMenuFunc) {

	h.MainMenuFunc = m

}

func (h *BotWithHandlers) registerMenuHandler(menu string, handler MenuHandlerFunc) {

	h.MenuHandlers[menu] = handler

}

func (h *BotWithHandlers) registerCommandHandler(command string, handler CommandHandlerFunc) {

	h.CommandHandlers[command] = handler

}

func (bot *BotWithHandlers) runStateHandler(ctx SessionContextInterface) {

	state, _ := ctx.Get()

	fn, ok := bot.StateHandlers[state]

	if ok {
		fn(ctx)

	}

}

func (h *BotWithHandlers) runMenuHandler(ctx MessageContextInterface){

	menu := ctx.Text()
	menuHandler, ok := h.MenuHandlers[menu]

	if ok {
		menuHandler(ctx)
	}

}

func (c *BotWithHandlers) runCommandHandler(ctx CommandContextInterface) {

	fn, ok := c.CommandHandlers[ctx.Command()]

	if !ok {
		return
	}

	fn(ctx)

}

//// all "Messages" fn
//func (h *BotWithHandlers) mainMessageProcessor(update tgbotapi.Update) {
//
//	ctx := h.NewMessageContext(update)
//
//	if ctx.IsCommand() {
//
//		go h.runCommandHandler(update.Message.Command(), update.Message.CommandArguments(), ctx)
//
//	} else {
//
//		ok := h.runMenuHandler(ctx)
//		if !ok {
//			h.runStateHandler(ctx)
//		}
//
//	}
//
//}

//// callback queries fn
//func (h *BotWithHandlers) mainCallbackQueryProcessor(update tgbotapi.Update) {
//
//	ctx := h.NewMessageContext(update)
//
//	callbackData := CallbackData(update.CallbackQuery.Data)
//	callbackHandler, ok := h.CallbackHandlers[callbackData.Action().Value()]
//
//	if ok {
//		go callbackHandler(callbackData, ctx)
//	} else if h.ErrorCallbackHandler != nil {
//
//		h.ErrorCallbackHandler(callbackData, ctx)
//	}

}
