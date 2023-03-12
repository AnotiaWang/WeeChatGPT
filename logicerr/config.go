package logicerr

import "errors"

var (
	ConfigFileNotFoundError   = errors.New("未找到配置文件 config.yml")
	OpenAISKNotSetError       = errors.New("请在 config.yml 中设置 OpenAI SecretKey (SK)")
	OpenAIModelNotSetError    = errors.New("请在 config.yml 中设置要使用的 OpenAI 模型")
	OpenAIEndpointNotSetError = errors.New("请在 config.yml 中设置 OpenAI API Endpoint")
)
