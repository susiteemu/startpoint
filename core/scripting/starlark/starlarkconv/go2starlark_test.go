package starlarkconv

import (
	"math/big"
	"testing"

	"go.starlark.net/starlark"
)

func TestConvertBigInt(t *testing.T) {
	var val *big.Int = big.NewInt(123)
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}
	wantedType := "int"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}
	wantedValue := val
	slInt := slVal.(starlark.Int).BigInt()
	if slInt.Cmp(wantedValue) != 0 {
		t.Errorf("Value got %d, wanted %d", slInt, wantedValue)
		return
	}
}

func TestConvertInt8(t *testing.T) {
	var val int8 = 123
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}

	wantedType := "int"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	wantedValue := val
	var ptr *int8 = new(int8)
	err = starlark.AsInt(slVal, ptr)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}
	if wantedValue != *ptr {
		t.Errorf("Value got %d, wanted %d", ptr, wantedValue)
		return
	}
}

func TestConvertInt16(t *testing.T) {
	var val int16 = 123
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}

	wantedType := "int"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	wantedValue := val
	var ptr *int16 = new(int16)
	err = starlark.AsInt(slVal, ptr)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}
	if wantedValue != *ptr {
		t.Errorf("Value got %d, wanted %d", ptr, wantedValue)
		return
	}
}

func TestConvertInt32(t *testing.T) {
	var val int32 = 123
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}

	wantedType := "int"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	wantedValue := val
	var ptr *int32 = new(int32)
	err = starlark.AsInt(slVal, ptr)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}
	if wantedValue != *ptr {
		t.Errorf("Value got %d, wanted %d", ptr, wantedValue)
		return
	}
}

func TestConvertInt64(t *testing.T) {
	var val int64 = 123
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}

	wantedType := "int"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	wantedValue := val
	var ptr *int64 = new(int64)
	err = starlark.AsInt(slVal, ptr)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}
	if wantedValue != *ptr {
		t.Errorf("Value got %d, wanted %d", ptr, wantedValue)
		return
	}
}

func TestConvertInt(t *testing.T) {
	var val int = 123
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}

	wantedType := "int"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	wantedValue := val
	var ptr *int = new(int)
	err = starlark.AsInt(slVal, ptr)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}
	if wantedValue != *ptr {
		t.Errorf("Value got %d, wanted %d", ptr, wantedValue)
		return
	}
}

func TestConvertFloat32(t *testing.T) {
	var val float32 = 123.456
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}

	wantedType := "float"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	wantedValue := float64(val)
	got, ok := starlark.AsFloat(slVal)
	if !ok {
		t.Errorf("Expected ok")
		return
	}
	if wantedValue != got {
		t.Errorf("Value got %f, wanted %f", got, wantedValue)
		return
	}
}

func TestConvertFloat64(t *testing.T) {
	var val float64 = 123.456
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}

	wantedType := "float"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	wantedValue := val
	got, ok := starlark.AsFloat(slVal)
	if !ok {
		t.Errorf("Expected ok")
		return
	}
	if wantedValue != got {
		t.Errorf("Value got %f, wanted %f", got, wantedValue)
		return
	}
}

func TestConvertBool(t *testing.T) {
	var val bool = true
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}

	wantedType := "bool"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	wantedValue := starlark.True
	got, ok := slVal.(starlark.Bool)
	if !ok {
		t.Errorf("Expected ok")
		return
	}

	equal, _ := starlark.Equal(got, wantedValue)
	if !equal {
		t.Errorf("Value got %v, wanted %v", got, wantedValue)
		return
	}
}

func TestConvertString(t *testing.T) {
	val := "Lorem ipsum"
	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}
	wantedType := "string"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}
	wantedValue := val
	if slVal.(starlark.String).GoString() != wantedValue {
		t.Errorf("Value got %s, wanted %s", slVal.String(), wantedValue)
		return
	}
}

func TestMap(t *testing.T) {
	val := make(map[string]interface{})
	val["one"] = 123
	val["two"] = "Jane"
	val["three"] = map[string][]string{
		"three-one": {"Matt", "Jake", "Ellison"},
	}

	slVal, err := Convert(val)
	if err != nil {
		t.Errorf("Error not expected %v", err)
		return
	}
	wantedType := "dict"
	if slVal.Type() != wantedType {
		t.Errorf("Type got %s, wanted %s", slVal.Type(), wantedType)
		return
	}

	slDict := slVal.(*starlark.Dict)
	for _, k := range slDict.Keys() {
		slDictVal, _, _ := slDict.Get(k)
		if k.(starlark.String).GoString() == "one" {
			nr := slDictVal.(starlark.Int).BigInt()
			if nr.Cmp(big.NewInt(123)) != 0 {
				t.Errorf("Value got %s, wanted %s", slDictVal.String(), "123")
				return
			}
		} else if k.(starlark.String).GoString() == "two" {
			str := slDictVal.(starlark.String)
			if str.GoString() != "Jane" {
				t.Errorf("Value got %s, wanted %s", slDictVal.String(), "Jane")
				return
			}
		} else if k.(starlark.String).GoString() == "three" {
			tkDictVal, _ := slDictVal.(*starlark.Dict)
			for _, tk := range tkDictVal.Keys() {
				if tk.(starlark.String).GoString() == "three-one" {
					tkVal, _, _ := tkDictVal.Get(tk)
					arr := tkVal.(*starlark.List)
					for i := 0; i < arr.Len(); i++ {
						arrEl := arr.Index(i).(starlark.String).GoString()
						if i == 0 && arrEl != "Matt" {
							t.Errorf("Value got %s, wanted %s", slDictVal.String(), "Matt")
							return
						}
						if i == 1 && arrEl != "Jake" {
							t.Errorf("Value got %s, wanted %s", slDictVal.String(), "Jake")
							return
						}
						if i == 2 && arrEl != "Ellison" {
							t.Errorf("Value got %s, wanted %s", slDictVal.String(), "Ellison")
							return
						}
					}
				}

			}
		}
	}
}
