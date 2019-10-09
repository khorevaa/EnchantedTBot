package HBot

import (
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

		ctx := NewCallbackContext(callbackQuery)

		callbackHandler, ok := bot.CallbackHandlers[ctx.CallbackData().Action().Value()]

		if ok {
			callbackHandler(ctx)
		}
	})

	bot.Server.HandleMessage("*", func(message *tbot.Message) {

		switch {

		case isCommand(message):

			ctx := NewCommandContext(message)
			bot.runCommandHandler(ctx)

		case bot.isMenu(message):

			ctx := NewMessageContext(message)
			bot.runMenuHandler(ctx)

		case bot.checkState(message):

			ctx := NewStateContext(message, bot.Sessions)
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

	if len(bot.Sessions.Get(chatID)) == 0 {
		return false
	}

	return true
}

// IsCommand returns true if message starts with a "bot_command" entity.
func isCommand(m *tbot.Message) bool {
	if m.Entities == nil || len(m.Entities) == 0 {
		return false
	}

	entity := (m.Entities)[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}

// Command checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is removed. Use
// CommandWithAt() if you do not want that.
func Command(m *tbot.Message) string {
	command := CommandWithAt(m)

	if i := strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}

	return command
}

// CommandWithAt checks if the message was a command and if it was, returns the
// command. If the Message was not a command, it returns an empty string.
//
// If the command contains the at name syntax, it is not removed. Use Command()
// if you want that.
func CommandWithAt(m *tbot.Message) string {
	if !isCommand(m) {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := (m.Entities)[0]
	return m.Text[1:entity.Length]
}

// CommandArguments checks if the message was a command and if it was,
// returns all text after the command name. If the Message was not a
// command, it returns an empty string.
//
// Note: The first character after the command name is omitted:
// - "/foo bar baz" yields "bar baz", not " bar baz"
// - "/foo-bar baz" yields "bar baz", too
// Even though the latter is not a command conforming to the spec, the API
// marks "/foo" as command entity.
func CommandArguments(m *tbot.Message) string {
	if !isCommand(m) {
		return ""
	}

	// IsCommand() checks that the message begins with a bot_command entity
	entity := (m.Entities)[0]
	if len(m.Text) == entity.Length {
		return "" // The command makes up the whole message
	}

	return m.Text[entity.Length+1:]
}
