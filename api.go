package HBot

import (
	"github.com/khorevaa/EnchantedTBot/builders/menu"
	"github.com/khorevaa/EnchantedTBot/builders/menuBuilder"
	"github.com/khorevaa/EnchantedTBot/types"
)

func (c *BotWithHandlers) CallbackHandler(callback types.CallbackActionInterface, h types.CallbackHandlerFunc) {

	c.registerCallbackHandler(callback, h)
}

func (c *BotWithHandlers) StateHandler(state string, h types.StateHandlerFunc) {
	c.registerStateHandler(state, h)
}

func (c *BotWithHandlers) CommandHandler(command string, h types.CommandHandlerFunc) {
	c.registerCommandHandler(command, h)
}

func (c *BotWithHandlers) RegisterChatStateMachine(csm types.SessionStorage) {
	c.Sessions = csm
}

func (h *BotWithHandlers) GetChatStateMachine() types.SessionStorage {

	return h.Sessions

}

func (h *BotWithHandlers) NewMenuWithOptions(opts ...menu.Option) *menu.MenuBuilder {

	h.Menu = &menu.Menu{}

	builder := &menu.MenuBuilder{}
	builder.Option(menu.WithMainMenu(h.Menu))
	builder.Option(opts...)

	return builder

}

func (h *BotWithHandlers) NewMenu() *menu.MenuBuilder {

	return h.NewMenuWithOptions(
		menu.ResizeKeyboard(true),
		menu.SelectiveKeyboard(true),
	)
}
