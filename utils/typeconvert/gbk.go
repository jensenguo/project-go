package typeconvert

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GBKToUTF8 将 GBK 数据转为 UTF-8 数据。
func GBKToUTF8(data []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SGBKToUTF8 将 GBK 字符串转为 UTF-8 字符串。
func SGBKToUTF8(s string) (string, error) {
	data, err := GBKToUTF8([]byte(s))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
