package conv

import "fmt"

func AssertAndConvert[T any](data map[string]interface{}, fieldName string) (T, error) {
	var result T
	_, has := data[fieldName]
	if !has {
		return result, fmt.Errorf("Request is missing %s: cannot extract value", fieldName)
	}
	result, ok := data[fieldName].(T)
	if !ok {
		return result, fmt.Errorf("Could not cast value %v for field %s to type %T", data[fieldName], fieldName, result)
	}
	return result, nil
}

func ConvertMapOfInterfaceToString(input interface{}) (map[string]interface{}, bool) {
	asMapInterface, isMapInterface := input.(map[interface{}]interface{})
	asMapString := map[string]interface{}{}
	if isMapInterface {
		for k, v := range asMapInterface {
			asMapString[fmt.Sprintf("%v", k)] = v
		}
	}
	return asMapString, isMapInterface

}
