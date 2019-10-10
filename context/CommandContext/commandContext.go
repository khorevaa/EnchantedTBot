package CommandContext

import (
	"github.com/khorevaa/EnchantedTBot/context/internal/baseContext"
	"github.com/khorevaa/EnchantedTBot/types"
	"github.com/yanzay/tbot/v2"
	"strconv"
	"strings"
)

type Option func(ctx *CommandContext) Option

var _ types.CommandContextInterface = (*CommandContext)(nil)

type CommandContext struct {
	baseContext.Context

	command   string
	args      []string
	message   *tbot.Message
	messageID int
}

func (ctx *CommandContext) Options(opts ...Option) (previous Option) {

	for _, opt := range opts {
		previous = opt(ctx)
	}
	return previous
}

func WithMessage(message *tbot.Message) Option {

	return func(ctx *CommandContext) Option {

		p := ctx.Message()

		chatID, _ := strconv.ParseInt(message.Chat.ID, 10, 64)
		userID := int64(message.From.ID)
		ctx.Options(WithChatID(chatID), WithUserID(userID))
		ctx.messageID = message.MessageID
		ctx.command = command(message)
		ctx.args = commandArguments(message)

		return WithMessage(p)
	}
}

func WithChatID(chatID int64) Option {

	return func(ctx *CommandContext) Option {
		p := ctx.ChatID()
		ctx.Context.Options(baseContext.WithChatID(chatID))
		return WithChatID(p)
	}

}

func WithUserID(userID int64) Option {

	return func(ctx *CommandContext) Option {
		p := ctx.UserID()
		ctx.Context.Options(baseContext.WithUserID(userID))
		return WithUserID(p)
	}

}

func WithLocal(local string) Option {

	return func(ctx *CommandContext) Option {
		p := ctx.GetLocal()
		ctx.Context.Options(baseContext.WithLocal(local))
		return WithLocal(p)
	}

}

func WithI18n(i18n types.TranslateInterface) Option {

	return func(ctx *CommandContext) Option {
		p := ctx.I18n()
		ctx.Context.Options(baseContext.WithI18n(i18n))
		return WithI18n(p)
	}

}

func WithMenu(fn types.TypeMainMenuFunc) Option {

	return func(ctx *CommandContext) Option {
		p := ctx.GetMenuMarkup()
		ctx.Context.Options(baseContext.WithMenu(fn))
		return WithMenu(p)
	}
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

// IsCommand returns true if message starts with a "bot_command" entity.
func isCommand(m *tbot.Message) bool {
	if m.Entities == nil || len(m.Entities) == 0 {
		return false
	}

	entity := (m.Entities)[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}

func command(m *tbot.Message) string {
	commandName := commandWithAt(m)

	if i := strings.Index(commandName, "@"); i != -1 {
		commandName = commandName[:i]
	}
	return commandName
}

func commandWithAt(m *tbot.Message) string {
	if !isCommand(m) {
		return ""
	}

	entity := (m.Entities)[0]
	return m.Text[1:entity.Length]
}

func commandArguments(m *tbot.Message) (args []string) {
	if !isCommand(m) {
		return
	}

	entity := (m.Entities)[0]
	if len(m.Text) == entity.Length {
		return
	}

	return strings.Split(m.Text[entity.Length+1:], " ")

}
