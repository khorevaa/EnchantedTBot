package csm

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type ChatStateMachineInterface interface {
	UpdateState(update tgbotapi.Update, state string) map[string]interface{}
	GetState(update tgbotapi.Update) string
	UpdateData(update tgbotapi.Update, data map[string]interface{}) map[string]interface{}
	GetData(update tgbotapi.Update) map[string]interface{}
	ClearData(update tgbotapi.Update)
}
