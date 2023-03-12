package logicerr

import "errors"

var (
	ChatCompletionFailedError = errors.New("请求接口失败")
	DecodeJSONFailedError     = errors.New("解析 JSON 失败")
)
