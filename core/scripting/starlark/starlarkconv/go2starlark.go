package starlarkconv

import (
	"errors"
	"go.starlark.net/starlark"
	"reflect"
)

var mustNotBeNil = errors.New("value must not be nil")

func Convert(v interface{}) (starlark.Value, error) {
	var converters = []interface{}{
		ConvertDict,
		ConvertString,
		ConvertBool,
		ConvertFloat,
		ConvertInt,
	}
	var converted starlark.Value
	for _, converter := range converters {
		converterFunc := converter.(func(v interface{}) (starlark.Value, bool, error))
		res, accept, err := converterFunc(v)
		if err != nil {
			return nil, err
		} else if accept {
			converted = res
			break
		}
	}
	return converted, nil
}

func typeOf(v interface{}) string {
	t := reflect.TypeOf(v)
	if t == nil {
		return ""
	}
	return t.String()
}

func ConvertInt(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	if typeOf(v) != "int" {
		return nil, false, nil
	}
	return starlark.MakeInt(v.(int)), true, nil
}

func ConvertFloat(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	if typeOf(v) != "float64" {
		return nil, false, nil
	}
	return starlark.Float(v.(float64)), true, nil
}

func ConvertString(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	if typeOf(v) != "string" {
		return nil, false, nil
	}
	return starlark.String(v.(string)), true, nil
}

func ConvertBool(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	if typeOf(v) != "bool" {
		return nil, false, nil
	}
	if v.(bool) {
		return starlark.True, true, nil
	}
	return starlark.False, true, nil
}

func ConvertDict(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Map {
		return nil, false, nil
	}
	d := starlark.Dict{}
	for _, key := range val.MapKeys() {
		mapValue := val.MapIndex(key)
		convertedKey, err := Convert(key.Interface())
		if err != nil {
			return nil, true, err
		}
		convertedVal, err := Convert(mapValue.Interface())
		if err != nil {
			return nil, true, err
		}
		err = d.SetKey(convertedKey, convertedVal)
		if err != nil {
			return nil, true, err
		}
	}
	return &d, true, nil
}
