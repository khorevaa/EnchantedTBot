package csm

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

type RedisCsm struct {
	*redis.Client
}

var _ ChatStateMachineInterface = (*RedisCsm)(nil)

// base set reimplementation with err handling
func (csm *RedisCsm) Set(key string, value string) {
	err := csm.Client.Set(key, value, time.Duration(DefaultExpirationTime)).Err()
	if err != nil {
		panic(err)
	}
}

// get data from key
// use `create=true` case you want to create new bucket with `key: {"__state__": INITIAL_VALUE}`
// if it was not found by `key`
func (csm *RedisCsm) Get(key string, create bool) map[string]interface{} {
	val, err := csm.Client.Get(key).Result()

	var result map[string]interface{}

	if err == redis.Nil {
		if create {
			result = map[string]interface{}{
				"__state__": InitialState,
			}
			csm.Set(key,
				fmt.Sprintf(`{"__state__": "%s"}`, InitialState)) // noqa

		} else {
			return result
		}
	} else if err != nil {
		panic(err)
	} else {
		err := json.Unmarshal([]byte(val), &result)
		if err != nil {
			panic(err)
		}
	}

	return result
}

// will create new record if can't get one with particular key
func (csm *RedisCsm) GetOrCreate(key string) map[string]interface{} {
	return csm.Get(key, true)
}

// get user identifier by most usable telegram update
// TODO: make all updates identifier
func (csm *RedisCsm) Key(update tgbotapi.Update) string {
	var FID int
	var CHID int64

	switch {
	case update.Message != nil:
		FID = update.Message.From.ID
		CHID = update.Message.Chat.ID
	case update.CallbackQuery != nil:
		FID = update.CallbackQuery.From.ID
		CHID = update.CallbackQuery.Message.Chat.ID
	case update.ChannelPost != nil:
		CHID = update.ChannelPost.Chat.ID
		FID = int(CHID)
	case update.ChosenInlineResult != nil:
		FID = update.ChosenInlineResult.From.ID
		CHID = int64(CHID)
	case update.EditedMessage != nil:
		FID = update.EditedMessage.From.ID
		CHID = update.EditedMessage.Chat.ID
	}

	return fmt.Sprintf(RecordPrefix+":%d:%d", FID, CHID)
}

// update data bucket
func (csm *RedisCsm) UpdateData(update tgbotapi.Update, data map[string]interface{}) map[string]interface{} {

	key := csm.Key(update)
	oldData := csm.GetOrCreate(key)

	for k, v := range data {
		oldData[k] = v
	}

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	csm.Set(key, string(body))
	return oldData
}

// update data bucket
func (csm *RedisCsm) ClearData(update tgbotapi.Update) {

	key := csm.Key(update)
	oldData := csm.GetOrCreate(key)

	var data = make(map[string]interface{}, 1)

	data["__state__"] = oldData["__state__"]

	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	csm.Set(key, string(body))

}

// sets new state
// TODO: discuss if __state__ key must be array of database and renamed to __states__
func (csm *RedisCsm) UpdateState(update tgbotapi.Update, state string) map[string]interface{} {
	key := csm.Key(update)
	oldData := csm.GetOrCreate(key)

	oldData["__state__"] = state

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	csm.Set(key, string(body))
	return oldData
}

func (csm *RedisCsm) GetState(update tgbotapi.Update) string {
	/*
		GetCurrentState used to get user's particular state
		:param: update - telegramBotUpdate
		:return: string state if declared
	*/
	key := csm.Key(update)
	oldData := csm.GetOrCreate(key)
	log.Println(oldData)
	return oldData["__state__"].(string)
}

// get user's saved data if no user found - will create new bucket
func (csm *RedisCsm) GetData(update tgbotapi.Update) map[string]interface{} {
	key := csm.Key(update)
	data := csm.GetOrCreate(key)
	return data
}

// remove all stored data and database (flush db)
func (csm *RedisCsm) Flush() {
	csm.FlushDB()
}
