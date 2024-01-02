package goconv

import (
	"errors"
	"fmt"
	"go.starlark.net/starlark"
)

var mustNotBeNil = errors.New("value must not be nil")

func ConvertValue(v starlark.Value) (interface{}, error) {
	var converters = []interface{}{
		ConvertDict,
		ConvertList,
		ConvertString,
		ConvertBool,
		ConvertNoneType,
		ConvertFloat,
		ConvertInt,
		ConvertBytes,
	}

	var converted interface{}
	for _, converter := range converters {
		converterFunc := converter.(func(v starlark.Value) (any, bool, error))
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

func ConvertDict(v starlark.Value) (interface{}, bool, error) {
	if v == nil {
		return map[string]interface{}{}, false, mustNotBeNil
	}
	if v.Type() != "dict" {
		return map[string]interface{}{}, false, nil
	}
	dict := v.(*starlark.Dict)
	tuples := dict.Items()
	goDict := make(map[string]interface{})
	for _, t := range tuples {
		t1 := t[0]
		t2 := t[1]
		converted, err := ConvertValue(t2)
		if err != nil {
			return map[string]interface{}{}, true, err
		}
		goDict[t1.String()] = converted
	}
	return goDict, true, nil
}

func ConvertString(v starlark.Value) (interface{}, bool, error) {
	if v == nil {
		return "", false, mustNotBeNil
	}
	if v.Type() != "string" {
		return "", false, nil
	}
	str := v.(starlark.String)
	return str.GoString(), true, nil
}

func ConvertNoneType(v starlark.Value) (interface{}, bool, error) {
	if v == nil {
		return "", false, mustNotBeNil
	}
	if v.Type() != "NoneType" {
		return "", false, nil
	}
	// TODO figure out proper way to convert this
	return "None", true, nil
}

func ConvertBool(v starlark.Value) (interface{}, bool, error) {
	if v == nil {
		return false, false, mustNotBeNil
	}
	if v.Type() != "bool" {
		return false, false, nil
	}
	equal, err := starlark.Equal(v, starlark.True)
	if err != nil {
		return false, true, err
	}
	return equal, true, nil
}

func ConvertFloat(v starlark.Value) (interface{}, bool, error) {
	if v == nil {
		return 0, false, mustNotBeNil
	}
	if v.Type() != "float" {
		return 0, false, nil
	}
	f, ok := starlark.AsFloat(v)
	if !ok {
		return 0, true, errors.New(fmt.Sprintf("failed to convert '%v' to float64", v))
	}
	return f, true, nil
}

func ConvertInt(v starlark.Value) (interface{}, bool, error) {
	if v == nil {
		return 0, false, mustNotBeNil
	}
	if v.Type() != "int" {
		return 0, false, nil
	}
	r, ok := v.(starlark.Int)
	if !ok {
		return 0, true, errors.New("could not get int")
	}
	return r.BigInt(), true, nil
}

func ConvertBytes(v starlark.Value) (interface{}, bool, error) {
	if v == nil {
		return "", false, mustNotBeNil
	}
	if v.Type() != "bytes" {
		return "", false, nil
	}
	bytes := v.(starlark.Bytes)
	return string(bytes), true, nil
}

func ConvertList(v starlark.Value) (interface{}, bool, error) {
	if v == nil {
		return []interface{}{}, false, mustNotBeNil
	}
	if v.Type() != "list" {
		return []interface{}{}, false, nil
	}
	var convertedItems []interface{}
	list := v.(*starlark.List)
	for i := 0; i < list.Len(); i++ {
		item := list.Index(i)
		converted, err := ConvertValue(item)
		if err != nil {
			return nil, true, err
		}
		convertedItems = append(convertedItems, converted)
	}
	return convertedItems, true, nil
}
