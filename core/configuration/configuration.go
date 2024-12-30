package configuration

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Configuration struct {
	requestOptions map[string]interface{}
}

func New() *Configuration {
	return &Configuration{}
}

func NewWithRequestOptions(options map[string]interface{}) *Configuration {
	return &Configuration{
		requestOptions: options,
	}
}

// TODO: rename this func
func (c *Configuration) GetStringOrDefault(key ...string) string {
	for _, k := range key {
		value, has := c.requestOptions[k]
		if has {
			return value.(string)
		}
		value = viper.GetString(k)
		if value != "" {
			return value.(string)
		}
	}
	return ""
}

func (c *Configuration) GetString(key string) (string, bool) {
	value, has := c.requestOptions[key]
	if has {
		return value.(string), true
	}
	if viper.IsSet(key) {
		return viper.GetString(key), true
	}
	return "", false
}

func (c *Configuration) GetStringSlice(key string) ([]string, bool) {
	value, has := c.requestOptions[key]
	log.Debug().Msgf("Key %s has %v value from request options %v", key, has, value)
	if has {
		strSlice := []string{}
		for _, v := range value.([]interface{}) {
			strSlice = append(strSlice, v.(string))
		}
		return strSlice, true
	}
	if viper.IsSet(key) {
		return viper.GetStringSlice(key), true
	}
	return []string{}, false
}

func (c *Configuration) Get(key string) (interface{}, bool) {
	value, has := c.requestOptions[key]
	if has {
		return value, true
	}
	if viper.IsSet(key) {
		return viper.Get(key), true
	}
	return nil, false
}

func (c *Configuration) GetSliceMapString(key string) ([]map[string]string, bool) {
	var cfgValue []interface{}
	value, has := c.requestOptions[key]
	// TODO: convert
	log.Debug().Msgf("Key %s has %v value from request options %v", key, has, value)
	if has {
		cfgValue, _ = value.([]interface{})
	} else {
		if viper.IsSet(key) {
			c := viper.Get(key)
			cfgValue, _ = c.([]interface{})
		}
	}

	if cfgValue != nil {
		result := make([]map[string]string, 0, 0)
		for _, s := range cfgValue {
			asMap, ok := s.(map[string]interface{})
			asMapString := make(map[string]string)
			if ok {
				for k, v := range asMap {
					asMapString[k] = v.(string)
				}
				result = append(result, asMapString)
			}
		}
		log.Debug().Msgf("Key %s has result %v", key, result)
		return result, true
	}

	return nil, false
}

func (c *Configuration) GetInt(key string) (int, bool) {
	value, has := c.requestOptions[key]
	if has {
		return value.(int), true
	}
	if viper.IsSet(key) {
		return viper.GetInt(key), true
	}
	return -1, false
}

func (c *Configuration) GetBool(key ...string) bool {
	for _, k := range key {
		value, has := c.requestOptions[k]
		if has {
			return value.(bool)
		}
		if viper.IsSet(k) {
			return viper.GetBool(k)
		}
	}
	return false
}

func (c *Configuration) GetBoolWithDefault(key string, dflt bool) bool {
	value, has := c.requestOptions[key]
	if has {
		return value.(bool)
	}
	if viper.IsSet(key) {
		return viper.GetBool(key)
	}
	return dflt
}

func Flatten(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	if len(prefix) > 0 {
		prefix += "."
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			Flatten(prefix+k, child, dest)
		case []interface{}:
			genKey := prefix + k
			dest[genKey] = child
			//for i := 0; i < len(child); i++ {
			//genKey := fmt.Sprintf("%s%s.%s", prefix, k, strconv.Itoa(i))
			//dest[genKey] = child[i]
			//}
		default:
			dest[prefix+k] = v
		}
	}
}
