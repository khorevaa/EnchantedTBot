package HBot

import (
	"github.com/khorevaa/EnchantedTBot/builders/menuBuilder"
	"github.com/khorevaa/EnchantedTBot/csm"
	"github.com/khorevaa/EnchantedTBot/handlers/Session"
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
)

type BotWithHandlers struct {
	*tbot.Server

	StateHandlers    *session.SessionHandlers
	MenuHandlers     map[string]types.MenuHandlerFunc
	CallbackHandlers map[int8]types.CallbackHandlerFunc
	CommandHandlers  map[string]types.CommandHandlerFunc

	Sessions SessionStorage

	Menu *menu.Menu

	MainMenuFunc types.TypeMainMenuFunc

	ErrorCallbackHandler func(callback types.CallbackDataInterface, ctx types.HandlerContextInterface)
}

func (h *BotWithHandlers) InitHandlers() {

	h.StateHandlers = NewStateHandlers()
	h.MenuHandlers = make(map[string]types.MenuHandlerFunc)
	h.CallbackHandlers = make(map[int8]types.CallbackHandlerFunc)
	h.CommandHandlers = make(map[string]types.CommandHandlerFunc)

}

func (h *BotWithHandlers) registerCallbackHandler(callback types.CallbackActionInterface, handler types.CallbackHandlerFunc) {

	key := callback.Value()
	h.CallbackHandlers[key] = handler

}

func (h *BotWithHandlers) registerStateHandler(state string, handler types.StateHandlerFunc) {

	if h.GetChatStateMachine() == nil {
		h.RegisterChatStateMachine(&csm.InMenoryCsm{})
	}

	h.StateHandlers.Add(state, handler)

}

func (h *BotWithHandlers) RegisterMenuFunc(m types.TypeMainMenuFunc) {

	h.MainMenuFunc = m

}

func (h *BotWithHandlers) registerMenuHandler(menu string, handler types.MenuHandlerFunc) {

	h.MenuHandlers[menu] = handler

}

func (h *BotWithHandlers) registerCommandHandler(command string, handler types.CommandHandlerFunc) {

	h.CommandHandlers[command] = handler

}

func (bot *BotWithHandlers) runStateHandler(ctx types.SessionContextInterface) {

	state, _ := ctx.Get()

	bot.StateHandlers.Run(state, ctx)

}

func (h *BotWithHandlers) runMenuHandler(ctx types.MessageContextInterface) {

	menu := ctx.Text()
	menuHandler, ok := h.MenuHandlers[menu]

	if ok {
		menuHandler(ctx)
	}

}

func (c *BotWithHandlers) runCommandHandler(ctx types.CommandContextInterface) {

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
