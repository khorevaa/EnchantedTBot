package HBot

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
	"time"
)

type SessionStorage interface {
	Set(int64, string)
	Get(int64) (string, map[string]interface{})
	SetData(int64, map[string]interface{})
	GetData(int64) map[string]interface{}
	Reset(int64)
}

var _ SessionStorage = (*RedisCsm)(nil)
var _ SessionStorage = (*InMemoryStorage)(nil)

const (
	DefaultExpirationTime int64  = 0           // records won't expire
	RecordPrefix          string = "user_data" // prefix for redis records
	InitialState          string = "~"
)

type InMemoryStorage struct {
	sync.Mutex
	sessions map[int64]string
}

func NewSessionStorage() SessionStorage {
	return &InMemoryStorage{sessions: make(map[int64]string)}
}

func (ims *InMemoryStorage) Get(id int64) string {
	ims.Lock()
	path := ims.sessions[id]
	ims.Unlock()
	return path
}

func (ims *InMemoryStorage) Set(id int64, path string) {
	ims.Lock()
	ims.sessions[id] = path
	ims.Unlock()
}

func (ims *InMemoryStorage) Reset(id int64) {
	ims.Lock()
	delete(ims.sessions, id)
	ims.Unlock()
}

type RedisCsm struct {
	*redis.Client
}

// base set reimplementation with err handling
func (csm *RedisCsm) set(key string, value string) {
	err := csm.Client.Set(key, value, time.Duration(DefaultExpirationTime)).Err()
	if err != nil {
		panic(err)
	}
}

// get data from key
// use `create=true` case you want to create new bucket with `key: {"__state__": INITIAL_VALUE}`
// if it was not found by `key`
func (csm *RedisCsm) get(key string, create bool) map[string]interface{} {
	val, err := csm.Client.Get(key).Result()

	var result map[string]interface{}

	if err == redis.Nil {
		if create {
			result = map[string]interface{}{
				"__state__": InitialState,
			}
			csm.set(key,
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
func (csm *RedisCsm) getOrCreate(key string) map[string]interface{} {
	return csm.get(key, true)
}

// get user identifier by most usable telegram update
// TODO: make all updates identifier
func (csm *RedisCsm) key(key int64) string {

	return fmt.Sprintf(RecordPrefix+":%d", key)
}

// update data bucket
func (csm *RedisCsm) SetData(sessionID int64, data map[string]interface{}) {

	key := csm.key(sessionID)

	oldData := csm.getOrCreate(key)

	for k, v := range data {
		oldData[k] = v
	}

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	csm.set(key, string(body))
}

// update data bucket
func (csm *RedisCsm) Reset(sessionID int64) {

	key := csm.key(sessionID)
	oldData := csm.getOrCreate(key)

	var data = make(map[string]interface{}, 1)

	data["__state__"] = oldData["__state__"]

	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	csm.set(key, string(body))

}

// sets new state
// TODO: discuss if __state__ key must be array of database and renamed to __states__
func (csm *RedisCsm) UpdateState(sessionID int64, state string) map[string]interface{} {
	key := csm.key(sessionID)
	oldData := csm.getOrCreate(key)

	oldData["__state__"] = state

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	csm.set(key, string(body))
	return oldData
}

// TODO: discuss if __state__ key must be array of database and renamed to __states__
func (csm *RedisCsm) Set(sessionID int64, state string) {
	key := csm.key(sessionID)
	oldData := csm.getOrCreate(key)

	oldData["__state__"] = state

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	csm.set(key, string(body))

}

func (csm *RedisCsm) Get(sessionID int64) (string, map[string]interface{}) {

	key := csm.key(sessionID)
	oldData := csm.getOrCreate(key)
	log.Println(oldData)
	return oldData["__state__"].(string), oldData
}

// get user's saved data if no user found - will create new bucket
func (csm *RedisCsm) GetData(sessionID int64) map[string]interface{} {
	key := csm.key(sessionID)
	data := csm.getOrCreate(key)
	return data
}

// remove all stored data and database (flush db)
func (csm *RedisCsm) Flush() {
	csm.FlushDB()
}
