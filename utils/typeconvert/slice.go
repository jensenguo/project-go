// Package typeconvert 数据类型转换
package typeconvert

import (
	"strconv"
	"unsafe"
)

// SliceU642I64 []uint64 数组转 []int64 数组
func SliceU642I64(s []uint64) []int64 {
	if s == nil {
		return nil
	}
	d := *(*[]int64)(unsafe.Pointer(&s))
	return d
}

// SliceU642I []uint64 数组转 []int 数组
func SliceU642I(s []uint64) []int {
	if s == nil {
		return nil
	}
	d := *(*[]int)(unsafe.Pointer(&s))
	return d
}

// SliceI642U64 []int64 数组转 []uint64 数组
func SliceI642U64(s []int64) []uint64 {
	if s == nil {
		return nil
	}
	d := *(*[]uint64)(unsafe.Pointer(&s))
	return d
}

// SliceU32ToI64 []int64 转 []uint32。
func SliceU32ToI64(s []uint32) []int64 {
	if s == nil {
		return nil
	}
	d := make([]int64, 0, len(s))
	for _, i32 := range s {
		d = append(d, int64(i32))
	}
	return d
}

// SliceI642U32 []int64 数组转 []uint32 数组
func SliceI642U32(s []int64) []uint32 {
	if s == nil {
		return nil
	}
	d := make([]uint32, 0, len(s))
	for _, i64 := range s {
		d = append(d, uint32(i64))
	}
	return d
}

// SliceI2I64 []int 数组转 []int64 数组
func SliceI2I64(s []int) []int64 {
	if s == nil {
		return nil
	}
	d := *(*[]int64)(unsafe.Pointer(&s))
	return d
}

// SliceI642I []int64 数组转 []int 数组
func SliceI642I(s []int64) []int {
	if s == nil {
		return nil
	}
	d := *(*[]int)(unsafe.Pointer(&s))
	return d
}

// SliceI642Interface []int64 数组转 []interface{} 数组
func SliceI642Interface(s []int64) []interface{} {
	if s == nil {
		return nil
	}
	d := make([]interface{}, 0, len(s))
	for _, i := range s {
		d = append(d, i)
	}
	return d
}

// SliceI82U8 convert []byte to []uint8
func SliceI82U8(s []byte) []uint8 {
	d := *(*[]uint8)(unsafe.Pointer(&s))
	return d
}

// SliceU82I8 convert []uint8 to []byte
func SliceU82I8(s []uint8) []byte {
	d := *(*[]byte)(unsafe.Pointer(&s))
	return d
}

// Bytes2Str []byte类型强转string类型，不进行数据拷贝
// 用户需要确保byte数组不被修改，否则字符串也会被修改
// 应用场景：json.Marshal时，返回的byte数组转为字符串
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// SliceI642SliceStr []int => []string
func SliceI642SliceStr(i64Sli []int64) []string {
	strSli := make([]string, len(i64Sli))
	for i, i64 := range i64Sli {
		strSli[i] = strconv.FormatInt(i64, 10)
	}
	return strSli
}

// SliceU322SliceStr []uint32 => []string
func SliceU322SliceStr(u32Sli []uint32) []string {
	strSli := make([]string, len(u32Sli))
	for i, u32 := range u32Sli {
		strSli[i] = strconv.FormatUint(uint64(u32), 10)
	}
	return strSli
}

// SliceI2SliceStr []int => []string
func SliceI2SliceStr(iSli []int) []string {
	strSli := make([]string, len(iSli))
	for i, i64 := range iSli {
		strSli[i] = strconv.Itoa(i64)
	}
	return strSli
}

// SliceStr2SliceI []string => []int, (不是数字或者空返回0，溢出返回极值)
func SliceStr2SliceI(strSli []string) []int {
	iSli := make([]int, len(strSli))
	for i, v := range strSli {
		iSli[i], _ = strconv.Atoi(v)
	}
	return iSli
}

// SliceStr2SliceI64 []string => []int64, 使用ParseInt(不是数字或者空返回0，溢出返回极值)
func SliceStr2SliceI64(strSli []string) []int64 {
	i64Sli := make([]int64, len(strSli))
	for i, v := range strSli {
		i64Sli[i], _ = strconv.ParseInt(v, 10, 64)
	}
	return i64Sli
}

// SliceStr2SliceInterface []string => []interface{}
func SliceStr2SliceInterface(strs []string) []interface{} {
	res := make([]interface{}, 0, len(strs))
	for _, str := range strs {
		res = append(res, str)
	}
	return res
}

// SliceInt32ToInt64 []int32 转 []int64。
func SliceInt32ToInt64(s []int32) []int64 {
	if s == nil {
		return nil
	}
	d := make([]int64, 0, len(s))
	for _, i32 := range s {
		d = append(d, int64(i32))
	}
	return d
}

// SliceInt64ToInt32 []int64 转 []int32。
func SliceInt64ToInt32(s []int64) []int32 {
	if s == nil {
		return nil
	}
	d := make([]int32, 0, len(s))
	for _, i64 := range s {
		d = append(d, int32(i64))
	}
	return d
}
