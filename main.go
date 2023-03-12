package main

import (
	"ChatGPT_WeChat_Bot/handler"
	"ChatGPT_WeChat_Bot/logicerr"
	"ChatGPT_WeChat_Bot/model"
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"log"
)

func main() {
	config := &model.Config{}

	if err := config.Load(); err != nil {
		if err == logicerr.ConfigFileNotFoundError {
			err = model.InitConfigFile()
			if err != nil {
				log.Fatal("创建配置文件 config.yml 失败: " + err.Error())
			}
			log.Println("已创建配置文件 config.yml，请参考文档，填写必要的字段。")
			return
		}
		panic(err)
	} else if err := config.Validate(); err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, model.ConfigKey, config)

	var bot *openwechat.Bot

	if config.WeChat.DesktopMode != nil && *config.WeChat.DesktopMode {
		bot = openwechat.DefaultBot(openwechat.Desktop)
	} else {
		bot = openwechat.DefaultBot()
	}

	reloadStorage := openwechat.NewFileHotReloadStorage("wechat_cache.json")
	defer reloadStorage.Close()

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx = context.WithValue(ctx, model.SelfKey, self)

	bot.MessageHandler = handler.Default(ctx)

	// 阻塞主 goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
