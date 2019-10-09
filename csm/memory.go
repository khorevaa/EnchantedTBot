package csm

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type InMenoryCsm struct {
	data map[string]string
}

const (
	DefaultExpirationTime int64  = 0           // records won't expire
	RecordPrefix          string = "user_data" // prefix for redis records
	InitialState          string = "~"
)

var _ ChatStateMachineInterface = (*InMenoryCsm)(nil)

// base set reimplementation with err handling
func (c *InMenoryCsm) Set(key string, value string) {

	c.data[key] = value

}

// get data from key
// use `create=true` case you want to create new bucket with `key: {"__state__": INITIAL_VALUE}`
// if it was not found by `key`
func (c *InMenoryCsm) Get(key string, create bool) map[string]interface{} {
	val, ok := c.data[key]

	var result map[string]interface{}

	if !ok {
		if create {
			result = map[string]interface{}{
				"__state__": InitialState,
			}
			c.Set(key,
				fmt.Sprintf(`{"__state__": "%s"}`, InitialState)) // noqa

		} else {
			return result
		}
	} else {
		err := json.Unmarshal([]byte(val), &result)
		if err != nil {
			panic(err)
		}
	}

	return result
}

// will create new record if can't get one with particular key
func (c *InMenoryCsm) GetOrCreate(key string) map[string]interface{} {
	return c.Get(key, true)
}

// get user identifier by most usable telegram update
// TODO: make all updates identifier
func (c *InMenoryCsm) Key(update tgbotapi.Update) string {
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
func (c *InMenoryCsm) UpdateData(update tgbotapi.Update, data map[string]interface{}) map[string]interface{} {

	key := c.Key(update)
	oldData := c.GetOrCreate(key)

	for k, v := range data {
		oldData[k] = v
	}

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	c.Set(key, string(body))
	return oldData
}

// update data bucket
func (c *InMenoryCsm) ClearData(update tgbotapi.Update) {

	key := c.Key(update)
	oldData := c.GetOrCreate(key)

	var data = make(map[string]interface{}, 1)

	data["__state__"] = oldData["__state__"]

	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	c.Set(key, string(body))

}

// sets new state
// TODO: discuss if __state__ key must be array of database and renamed to __states__
func (c *InMenoryCsm) UpdateState(update tgbotapi.Update, state string) map[string]interface{} {
	key := c.Key(update)
	oldData := c.GetOrCreate(key)

	oldData["__state__"] = state

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	c.Set(key, string(body))
	return oldData
}

func (c *InMenoryCsm) GetState(update tgbotapi.Update) string {
	/*
		GetCurrentState used to get user's particular state
		:param: update - telegramBotUpdate
		:return: string state if declared
	*/
	key := c.Key(update)
	oldData := c.GetOrCreate(key)
	log.Println(oldData)
	return oldData["__state__"].(string)
}

// get user's saved data if no user found - will create new bucket
func (c *InMenoryCsm) GetData(update tgbotapi.Update) map[string]interface{} {
	key := c.Key(update)
	data := c.GetOrCreate(key)
	return data
}
