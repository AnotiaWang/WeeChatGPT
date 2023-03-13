package model

import (
	"ChatGPT_WeChat_Bot/logicerr"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type contextKey struct {
	name string
}

type Config struct {
	WeChat *WeChatConfig `yaml:"wechat"`
	OpenAI *OpenAIConfig `yaml:"openai"`
}

var ConfigKey = &contextKey{"config"}
var SelfKey = &contextKey{"self"}

type WeChatConfig struct {
	// 使用 PC 模式而不是网页版，或许可以解决部分新号无法使用的问题
	DesktopMode *bool `yaml:"desktopMode"`
}

type OpenAIConfig struct {
	Endpoint  string `yaml:"endpoint"`
	Model     string `yaml:"model"`
	SecretKey string `yaml:"secretKey"`
	Prefix    string `yaml:"prefix"`
	SystemMsg string `yaml:"systemMsg"`
}

func InitConfigFile() error {
	str := `# ChatGPT WeChat Bot 配置文件
wechat:
  # 是否使用 PC 模式而不是网页版，或许可以解决部分新号无法使用的问题
  desktopMode: false
openai:
  endpoint: https://api.openai.com/v1/chat/completions
  # 模型
  model: gpt-3.5-turbo
  # OpenAI SecretKey ( https://platform.openai.com/account/api-keys )
  secretKey: 
  # 前缀
  prefix: ChatGPT,
  # 设定身份，如：你是 ChatGPT，OpenAI 发布的语言模型
  systemMsg:
`
	err := os.WriteFile("./config.yml", []byte(str), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (c *Config) Load() error {
	yamlFile, err := os.ReadFile("./config.yml")
	if err != nil {
		if os.IsNotExist(err) {
			return logicerr.ConfigFileNotFoundError
		}
		fmt.Println(err.Error())
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func (c *Config) Validate() error {
	if c.OpenAI.SecretKey == "" {
		return logicerr.OpenAISKNotSetError
	} else if c.OpenAI.Model == "" {
		return logicerr.OpenAIModelNotSetError
	} else if c.OpenAI.Endpoint == "" {
		return logicerr.OpenAIEndpointNotSetError
	} else if c.OpenAI.Prefix == "" {
		c.OpenAI.Prefix = "ChatGPT,"
	}
	return nil
}
