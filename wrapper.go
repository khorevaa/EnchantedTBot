package HBot

import (
	"github.com/khorevaa/EnchantedTBot/context"
	"github.com/khorevaa/EnchantedTBot/context/CommandContext"
	"github.com/yanzay/tbot/v2"
	"strconv"
	"strings"
)

/*
New creates new Server. Available options:
	WithWebook(url, addr string)
	WithHTTPClient(client *http.Client)
*/
func New(token string, options ...tbot.ServerOption) *BotWithHandlers {

	bot := &BotWithHandlers{
		Server: tbot.New(token, options...),
	}

	bot.Server.HandleCallback(func(callbackQuery *tbot.CallbackQuery) {

		ctx := context.NewCallbackContext(callbackQuery)

		callbackHandler, ok := bot.CallbackHandlers[ctx.CallbackData().Action().Value()]

		if ok {
			callbackHandler(ctx)
		}
	})

	bot.Server.HandleMessage("*", func(message *tbot.Message) {

		switch {

		case isCommand(message):

			ctx := context.NewCommandContext(message, CommandContext.WithMenu(bot.MainMenuFunc))
			bot.runCommandHandler(ctx)

		case bot.isMenu(message):

			ctx := context.NewMessageContext(message)
			bot.runMenuHandler(ctx)

		case bot.checkState(message):

			ctx := context.NewStateContext(message, bot.Sessions)
			bot.runStateHandler(ctx)

		default:
			// none
		}

	})

	return bot

}

func (bot *BotWithHandlers) isMenu(m *tbot.Message) bool {

	if len(m.Text) == 0 {
		return false
	}

	_, ok := bot.MenuHandlers[m.Text]

	return ok
}

func (bot *BotWithHandlers) checkState(m *tbot.Message) bool {

	chatID, _ := strconv.ParseInt(m.Chat.ID, 10, 64)

	state, _ := bot.Sessions.Get(chatID)

	return bot.StateHandlers.Contain(state)
}

// IsCommand returns true if message starts with a "bot_command" entity.
func isCommand(m *tbot.Message) bool {
	if m.Entities == nil || len(m.Entities) == 0 {
		return false
	}

	entity := (m.Entities)[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}
