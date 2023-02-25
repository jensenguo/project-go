package typeconvert

import (
	"strconv"
	"unsafe"
)

// Str2Bytes string类型强转[]byte，不进行数据拷贝，
// 用户需要确保byte数组不被修改，否则字符串也会被修改，
// 如果修改字符串字面量，则会导致panic，
// 应用场景：json.UnMarshal时，第一个参数可以使用json.Unmarshal(str2bytes(s), v);
func Str2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

//Str2Int32 string转int32
func Str2Int32(s string) (int32, error) {
	if len(s) == 0 {
		return 0, nil
	}
	res, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(res), nil
}

//Str2Int string转int
func Str2Int(s string) (int, error) {
	if len(s) == 0 {
		return 0, nil
	}
	res, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return res, nil
}

// Str2Uint32 string转uint32
func Str2Uint32(s string) (uint32, error) {
	if len(s) == 0 {
		return 0, nil
	}
	u, err := strconv.ParseUint(s, 10, 32)
	return uint32(u), err
}

// Str2Uint16 string转uint16
func Str2Uint16(s string) (uint16, error) {
	if len(s) == 0 {
		return 0, nil
	}
	u, err := strconv.ParseUint(s, 10, 16)
	return uint16(u), err
}

// Str2Int64WithDefault string转int64，支持默认参数
func Str2Int64WithDefault(v string, defaultValue int64) int64 {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}
