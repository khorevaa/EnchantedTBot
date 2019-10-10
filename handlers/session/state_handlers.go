package session

import (
	"github.com/khorevaa/EnchantedTBot/types"
	"sort"
	"strings"
)

type SessionHandlers struct {
	registeredStates []string
	handlers         map[string]types.StateHandlerFunc
}

func (sh *SessionHandlers) Add(state string, fn types.StateHandlerFunc) {

	sh.registeredStates = append(sh.registeredStates, state)
	sh.handlers[state] = fn

}

func (sh *SessionHandlers) Find(state string) int {

	return sort.SearchStrings(sh.registeredStates, state)
}

func (sh *SessionHandlers) Contain(state string) bool {

	return strings.EqualFold(sh.registeredStates[sh.Find(state)], state)

}

func (sh *SessionHandlers) Run(state string, ctx types.SessionContextInterface) {

	fn, ok := sh.handlers[state]

	if ok {
		fn(ctx)
	}

}
