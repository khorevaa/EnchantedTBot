package Redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/khorevaa/EnchantedTBot/types"
	"log"
	"time"
)

var _ types.SessionStorage = (*RedisStorage)(nil)

const (
	DefaultExpirationTime int64  = 0           // records won't expire
	RecordPrefix          string = "user_data" // prefix for redis records
	InitialState          string = "~"
)

type RedisStorage struct {
	*redis.Client
}

// base set reimplementation with err handling
func (csm *RedisStorage) set(key string, value string) {
	err := csm.Client.Set(key, value, time.Duration(DefaultExpirationTime)).Err()
	if err != nil {
		panic(err)
	}
}

// get data from key
// use `create=true` case you want to create new bucket with `key: {"__state__": INITIAL_VALUE}`
// if it was not found by `key`
func (csm *RedisStorage) get(key string, create bool) map[string]interface{} {
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
func (csm *RedisStorage) getOrCreate(key string) map[string]interface{} {
	return csm.get(key, true)
}

// get user identifier by most usable telegram update
// TODO: make all updates identifier
func (csm *RedisStorage) key(key int64) string {

	return fmt.Sprintf(RecordPrefix+":%d", key)
}

// update data bucket
func (csm *RedisStorage) SetData(sessionID int64, data map[string]interface{}) {

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
func (csm *RedisStorage) Reset(sessionID int64) {

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
func (csm *RedisStorage) UpdateState(sessionID int64, state string) map[string]interface{} {
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
func (csm *RedisStorage) Set(sessionID int64, state string) {
	key := csm.key(sessionID)
	oldData := csm.getOrCreate(key)

	oldData["__state__"] = state

	body, err := json.Marshal(oldData)
	if err != nil {
		panic(err)
	}

	csm.set(key, string(body))

}

func (csm *RedisStorage) Get(sessionID int64) (string, map[string]interface{}) {

	key := csm.key(sessionID)
	oldData := csm.getOrCreate(key)
	log.Println(oldData)
	return oldData["__state__"].(string), oldData
}

// get user's saved data if no user found - will create new bucket
func (csm *RedisStorage) GetData(sessionID int64) map[string]interface{} {
	key := csm.key(sessionID)
	data := csm.getOrCreate(key)
	return data
}

// remove all stored data and database (flush db)
func (csm *RedisStorage) Flush() {
	csm.FlushDB()
}
