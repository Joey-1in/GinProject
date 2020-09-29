package tool

import (
	"io"
	"encoding/json"
)

type JsonParse struct {

}

// 统一解析参数的方法
func Decode(io io.ReadCloser, v interface{}) error {
	return json.NewDecoder(io).Decode(v)
}