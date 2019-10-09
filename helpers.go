package HBot

import (
	"encoding/hex"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"prizmlive-bot/internal/liveBot/HBot/csm"
	"prizmlive-bot/internal/liveBot/HBot/i18n"
)

func NewContextWithUpdate(update tgbotapi.Update) HandlerContextInterface {

	return &baseContext{
		i18n:            i18n.Translotor(),
		localTranslator: i18n.Translotor(),
	}

}

// redis db setup
func NewRedisCsm(opts redis.Options) *csm.RedisCsm {
	//Client = redis.NewClient(&redis.Options{
	//	Addr:         "0.0.0.0:6379",
	//	DialTimeout:  10 * time.Second,
	//	ReadTimeout:  30 * time.Second,
	//	WriteTimeout: 30 * time.Second,
	//	PoolSize:     100,
	//	PoolTimeout:  30 * time.Second,
	//})

	return &csm.RedisCsm{
		redis.NewClient(&opts),
	}

}

func getChatID(update tgbotapi.Update) int64 {

	var chatID int64

	if isCallback(update) {
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		chatID = update.Message.Chat.ID
	}

	return chatID
}

func isCallback(update tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

func toHexString(d []byte) string {
	return hex.EncodeToString(d)
}

func fromHexString(d string) []byte {

	b, e := hex.DecodeString(d)

	if e != nil {
		log.Panic(e)
	}

	return b
}

func getCallbackUserID(callback *tgbotapi.CallbackQuery) int64 {

	switch {

	case callback.Message.From != nil && !callback.Message.From.IsBot:
		return int64(callback.Message.From.ID)
	case callback.From != nil && !callback.From.IsBot:
		return int64(callback.From.ID)
	case callback.Message.Chat != nil:
		return callback.Message.Chat.ID

	default:
		log.Fatal("Error get user id")
		return -1
	}
}

func getUserID(update tgbotapi.Update) int64 {

	var chatID int64

	switch {

	case update.CallbackQuery != nil:

		return getCallbackUserID(update.CallbackQuery)

	case update.Message != nil:

		chatID = int64(update.Message.From.ID)

	}

	return chatID
}
