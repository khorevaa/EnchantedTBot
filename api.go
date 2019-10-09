package HBot

func (c *BotWithHandlers) CallbackHandler(callback CallbackActionInterface, h CallbackHandlerFunc) {

	c.registerCallbackHandler(callback, h)
}

func (c *BotWithHandlers) StateHandler(state string, h StateHandlerFunc) {
	c.registerStateHandler(state, h)
}

func (c *BotWithHandlers) MenuHandler(menu string, h MenuHandlerFunc) {
	c.registerMenuHandler(menu, h)
}

func (c *BotWithHandlers) CommandHandler(command string, h CommandHandlerFunc) {
	c.registerCommandHandler(command, h)
}

func (c *BotWithHandlers) RegisterChatStateMachine(csm SessionStorage) {
	c.Sessions = csm
}

func (h *BotWithHandlers) GetChatStateMachine() SessionStorage {

	return h.Sessions

}
