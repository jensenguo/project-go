package typeconvert

import (
	"strconv"
)

// MapStrStr2Str 获取map m中key对应对值，支持默认参数
func MapStrStr2Str(m map[string]string, key, defaultValue string) string {
	value, ok := m[key]
	if !ok {
		return defaultValue
	}
	return value
}

// MapStrStr2I64 获取map m中key对应对值，并转换成int64，支持默认参数
func MapStrStr2I64(m map[string]string, key string, defaultValue int64) int64 {
	value, ok := m[key]
	if !ok {
		return defaultValue
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// MapStrInterface2MapStrInterface ...
func MapStrInterface2MapStrInterface(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	value, ok := m[key]
	if !ok {
		return nil
	}
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	return valueMap
}
