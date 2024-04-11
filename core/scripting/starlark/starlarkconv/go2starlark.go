package starlarkconv

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/rs/zerolog/log"
	"go.starlark.net/starlark"
)

var mustNotBeNil = errors.New("value must not be nil")

func Convert(v interface{}) (starlark.Value, error) {
	var converters = []interface{}{
		ConvertDict,
		ConvertString,
		ConvertBool,
		ConvertFloat,
		ConvertInt,
		ConvertBigInt,
		ConvertUint,
		ConvertArray,
	}
	var converted starlark.Value
	for _, converter := range converters {
		converterFunc := converter.(func(v interface{}) (starlark.Value, bool, error))
		res, accept, err := converterFunc(v)
		log.Debug().Msgf("Converter func: %v, %v, %v", res, accept, err)
		if err != nil {
			return nil, err
		} else if accept {
			converted = res
			break
		}
	}

	if converted == nil {
		return nil, errors.New(fmt.Sprintf("Could not find converter for value %v with type %T", v, v))
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

func ConvertBigInt(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	t := typeOf(v)
	if t != "*big.Int" {
		return nil, false, nil
	}
	return starlark.MakeBigInt(v.(*big.Int)), true, nil
}

func ConvertInt(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	t := typeOf(v)
	if t != "int" && t != "int8" && t != "int16" && t != "int32" && t != "int64" {
		return nil, false, nil
	}

	if t == "int64" {
		return starlark.MakeInt64(v.(int64)), true, nil
	}

	var asInt int
	if t == "int" {
		asInt = v.(int)
	} else if t == "int8" {
		asInt = int(v.(int8))
	} else if t == "int16" {
		asInt = int(v.(int16))
	} else {
		asInt = int(v.(int32))
	}

	return starlark.MakeInt(asInt), true, nil
}

func ConvertUint(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	t := typeOf(v)
	if t != "uint" && t != "uint8" && t != "uint16" && t != "uint32" && t != "uint64" {
		return nil, false, nil
	}

	if t == "uint64" {
		return starlark.MakeUint64(v.(uint64)), true, nil
	}

	var asUint uint
	if t == "uint" {
		asUint = v.(uint)
	} else if t == "uint8" {
		asUint = uint(v.(uint8))
	} else if t == "uint16" {
		asUint = uint(v.(uint16))
	} else {
		asUint = uint(v.(uint32))
	}

	return starlark.MakeUint(asUint), true, nil
}

func ConvertFloat(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	t := typeOf(v)
	if t != "float32" && t != "float64" {
		return nil, false, nil
	}

	var asFloat64 float64
	if t == "float32" {
		asFloat64 = float64(v.(float32))
	} else {
		asFloat64 = v.(float64)
	}
	return starlark.Float(asFloat64), true, nil
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
	log.Debug().Msgf("Convert dict %v", v)
	d := starlark.Dict{}
	for _, key := range val.MapKeys() {
		log.Debug().Msgf(">> Key %v", key)
		mapValue := val.MapIndex(key)
		convertedKey, err := Convert(key.Interface())
		if err != nil {
			return nil, true, err
		}
		log.Debug().Msgf("key=%s, convertedKey=%v", key, convertedKey)
		convertedVal, err := Convert(mapValue.Interface())
		if err != nil {
			return nil, true, err
		}
		err = d.SetKey(convertedKey, convertedVal)
		if err != nil {
			return nil, true, err
		}
	}
	log.Debug().Msgf("Starlark dict %v", &d)
	return &d, true, nil
}

func ConvertArray(v interface{}) (starlark.Value, bool, error) {
	if v == nil {
		return nil, false, mustNotBeNil
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		return nil, false, nil
	}
	l := starlark.List{}
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		converted, err := Convert(item.Interface())
		if err != nil {
			return nil, true, err
		}
		l.Append(converted)
	}
	return &l, true, nil
}
