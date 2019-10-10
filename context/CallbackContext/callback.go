package CallbackContext

import (
	"encoding/hex"
	"fmt"
	"github.com/khorevaa/EnchantedTBot/types"
	"log"
	"strings"
)

var _ types.CallbackDataInterface = (*CallbackData)(nil)
var _ types.CallbackActionInterface = (*CallbackAction)(nil)
var _ types.CallbackActionDataInterface = (*ActionData)(nil)

const ErrorCallback CallbackAction = 127

const (
	sep     = "?"
	sepData = "&"
	keyData = "__data__"
)

type CallbackAction int8

func (c CallbackAction) String() string {

	return toHexString([]byte{byte(c)})

}

func (c CallbackAction) Value() int8 {

	return int8(c)

}

type CallbackData string

func (d CallbackData) String() string {

	return string(d)

}

func (d CallbackData) Value() string {

	return string(d)

}

type ActionData string

func (a ActionData) String() string {

	return string(a)
}

func (a ActionData) Value() string {

	return string(a)
}

func (a ActionData) Map() map[string]string {

	args := strings.Split(a.String(), sepData)
	result := make(map[string]string, len(args))

	for _, str := range args {

		keyVal := strings.Split(str, "=")

		if len(keyVal) == 1 {
			result[keyData] += str + ";"
		} else {
			result[keyVal[0]] = keyVal[1]
		}
	}

	return result
}

func (a ActionData) FromSlice(args ...string) types.CallbackActionDataInterface {

	var tmpSlice []string

	if len(args)%2 == 1 {

		return ActionData(strings.Join(args, sepData))

	}

	for i := 0; i < len(args); {

		if i == len(args)-1 {
			break
		}

		key, val := args[i], args[i+1]

		if len(tmpSlice) == 0 {
			tmpSlice = make([]string, 0, len(args)/2)
		}

		tmpSlice = append(tmpSlice, fmt.Sprintf("%s=%s", key, val))

		i += 2
	}

	a = ActionData(strings.Join(tmpSlice, sepData))

	return a
}

func (a ActionData) FromMap(in map[string]string) types.CallbackActionDataInterface {

	var args []string

	for k, v := range in {

		args = append(args, fmt.Sprintf("%s=%s", k, v))

	}

	a = ActionData(strings.Join(args, sepData))
	return a

}

func getCallbackData(callbackString string) (action CallbackAction, route []byte, data ActionData) {

	if len(callbackString) == 0 {
		return ErrorCallback, route, data
	}

	a := strings.Split(callbackString, sep)

	hexRoute := a[0]

	if len(a) == 2 {
		data = ActionData(strings.TrimPrefix(callbackString, hexRoute+sep))
	}

	route = fromHexString(hexRoute)
	action = CallbackAction(route[len(route)-1])

	return
}

func (d CallbackData) Action() types.CallbackActionInterface {

	a, _, _ := getCallbackData(d.String())

	return a
}

func (d CallbackData) Data() types.CallbackActionDataInterface {

	_, _, data := getCallbackData(d.String())

	return data
}

func (d CallbackData) Back(data ...string) string {

	_, route, _ := getCallbackData(d.String())

	if len(route)-1 < 0 {
		return ""
	}

	back := toHexString(route[:len(route)-1])

	if len(data) > 0 {
		actionData := ActionData("").FromSlice(data...).String()
		back += HBot.sep + actionData
	}

	return back
}

func (d CallbackData) Next(next types.CallbackActionInterface, data ...string) string {

	_, route, _ := getCallbackData(d.String())

	route = append(route, byte(next.Value()))

	nextData := toHexString(route)

	if len(data) > 0 {
		actionData := ActionData("").FromSlice(data...).String()
		nextData += sep + actionData
	}

	return nextData
}

func (d CallbackData) WithData(data ...string) types.CallbackDataInterface {

	if len(data) > 0 {
		actionData := ActionData("").FromSlice(data...).String()
		return CallbackData(d.String() + sep + actionData)
	}

	return d
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
