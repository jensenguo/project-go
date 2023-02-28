// Package json json相关函数
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

// 错误码
var (
	ErrDataEmpty         = fmt.Errorf("data empty")          // 数据是空的
	ErrDataFormatInvalid = fmt.Errorf("data format invalid") // 数据格式无效

	jsonGo            = jsoniter.ConfigCompatibleWithStandardLibrary                            // json go 库默认实体
	jsonPBUnmarshaler = jsonpb.Unmarshaler{AllowUnknownFields: true}                            // jsonpb unmarshal
	jsonPBMarshaller  = jsonpb.Marshaler{EmitDefaults: true, OrigName: true, EnumsAsInts: true} // jsonpb marshal
)

// init 全局初始化
func init() {
	// jsonGo 进行注册
	extra.RegisterFuzzyDecoders()
}

// MarshalString 利用标准库json任意类型转字符串，失败返回""
func MarshalString(v interface{}) string {
	bf := bytes.NewBuffer(nil)
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(v); err != nil {
		return ""
	}
	// 删除多余的换行
	by := bf.Bytes()
	if len(by) > 0 && by[len(by)-1] == '\n' {
		by = by[:len(by)-1]
	}
	return string(by)
}

// GoUnmarshal 利用 json-go 对json进行解析
func GoUnmarshal(data []byte, v interface{}) error {
	return jsonGo.Unmarshal(data, v)
}

// GoUnmarshalFile 利用jsongo对json 文件进行解析
func GoUnmarshalFile(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return jsonGo.Unmarshal(data, v)
}

// Copy 拷贝
func Copy(s, d interface{}) error {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("json.Marshal(s)  fail err=[%v]", err)
	}
	if err = jsonGo.Unmarshal(b, d); err != nil {
		return fmt.Errorf("json.Unmarshal d fail, err=[%v]", err)
	}
	return nil
}

// Form json 转 url form格式
func Form(v interface{}) (string, error) {
	m := make(map[string]json.RawMessage)
	err := Copy(v, &m)
	if err != nil {
		return "", nil
	}
	uv := &url.Values{}
	for k, v := range m {
		s := string(v)
		if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
			var s1 string
			if err = json.Unmarshal(v, &s1); err == nil {
				s = s1
			}
		}
		uv.Add(k, s)
	}
	return uv.Encode(), nil
}

// PbUnmarshal jsonpb 扩展map，slice，json 解析
func PbUnmarshal(data []byte, message proto.Message) error {
	return jsonPBUnmarshaler.Unmarshal(bytes.NewReader(data), message)
}

// PbMarshal jsonpb 扩展map，slice 转json
func PbMarshal(message proto.Message) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	if err := jsonPBMarshaller.Marshal(w, message); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// PbMarshalString jsonpb 扩展map，slice 转json string
func PbMarshalString(message proto.Message) string {
	buf, err := PbMarshal(message)
	if err != nil {
		return ""
	}
	return string(buf)
}

// JsonpUnmarshal 利用jsongo对json进行解析，并去掉外层对jsonp封装
func JsonpUnmarshal(in []byte) ([]byte, error) {
	beg := 0
	for ; beg < len(in); beg++ {
		if in[beg] == '(' {
			beg++
			break
		}
	}
	end := len(in) - 1
	for ; end >= 0; end-- {
		if in[end] == ')' {
			end--
			break
		}
	}
	if beg > end {
		return nil, fmt.Errorf("jsonp format invalid in=[%s]", in)
	}
	out := in[beg : end+1]
	return out, nil
}

// GoUnmarshalJsonp 利用jsongo对json进行解析，并去掉外层对jsonp封装
func GoUnmarshalJsonp(jsonpData []byte, v interface{}) error {
	jsonData, err := JsonpUnmarshal(jsonpData)
	if err != nil {
		return err
	}
	return jsonGo.Unmarshal(jsonData, v)
}
