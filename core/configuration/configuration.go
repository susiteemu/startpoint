package configuration

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func GetStringOrDefault(key string) string {
	log.Debug().Msgf("Getting configuration value with key %s", key)
	value := viper.GetString(key)
	log.Debug().Msgf("Configuration value with key %s is %s", key, value)
	return value
}

func GetString(key string) (string, bool) {
	if viper.IsSet(key) {
		return viper.GetString(key), true
	}
	return "", false
}

func GetStringSlice(key string) ([]string, bool) {
	if viper.IsSet(key) {
		return viper.GetStringSlice(key), true
	}
	return []string{}, false
}

func Get(key string) (interface{}, bool) {
	if viper.IsSet(key) {
		return viper.Get(key), true
	}
	return nil, false
}

func GetSliceMapString(key string) ([]map[string]string, bool) {
	if viper.IsSet(key) {
		result := make([]map[string]string, 0, 0)
		c := viper.Get(key)
		asSlice, ok := c.([]interface{})
		if ok {
			for _, s := range asSlice {
				asMap, ok := s.(map[string]interface{})
				asMapString := make(map[string]string)
				if ok {
					for k, v := range asMap {
						asMapString[k] = v.(string)
					}
					result = append(result, asMapString)
				}
			}
			return result, true
		}

	}
	return nil, false
}

func GetInt(key string) (int, bool) {
	if viper.IsSet(key) {
		return viper.GetInt(key), true
	}
	return -1, false
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}
