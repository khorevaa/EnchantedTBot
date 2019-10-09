package HBot

import (
	"fmt"
	"strings"
)

var _ CallbackDataInterface = (*CallbackData)(nil)
var _ CallbackActionInterface = (*CallbackAction)(nil)
var _ CallbackActionDataInterface = (*ActionData)(nil)

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

func (a ActionData) FromSlice(args ...string) CallbackActionDataInterface {

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

func (a ActionData) FromMap(in map[string]string) CallbackActionDataInterface {

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

func (d CallbackData) Action() CallbackActionInterface {

	a, _, _ := getCallbackData(d.String())

	return a
}

func (d CallbackData) Data() CallbackActionDataInterface {

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
		back += sep + actionData
	}

	return back
}

func (d CallbackData) Next(next CallbackActionInterface, data ...string) string {

	_, route, _ := getCallbackData(d.String())

	route = append(route, byte(next.Value()))

	nextData := toHexString(route)

	if len(data) > 0 {
		actionData := ActionData("").FromSlice(data...).String()
		nextData += sep + actionData
	}

	return nextData
}

func (d CallbackData) WithData(data ...string) CallbackDataInterface {

	if len(data) > 0 {
		actionData := ActionData("").FromSlice(data...).String()
		return CallbackData(d.String() + sep + actionData)
	}

	return d
}
