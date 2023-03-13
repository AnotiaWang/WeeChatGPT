# WeeChatGPT

将 ChatGPT 集成到微信个人号。基于 [openwechat](https://github.com/eatmoreapple/openwechat)。

## 运行方法

1. 前往 [Releases](https://github.com/AnotiaWang/WeeChatGPT/releases/latest) ，根据你使用的平台，下载最新版本的 WeeChatGPT。
2. 先运行一次 WeeChatGPT，程序会生成配置文件 `config.yml`，然后根据文件中的提示，填写配置。
3. 再次运行 WeeChatGPT，会显示微信的登录二维码链接，在浏览器中打开它，扫码登录微信即可。

### 提示

- 程序支持热登录，如果两次登录之间的间隔较短，可以自动登录。
- 建议将群组添加到通讯录，否则机器人可能无法获取到群组信息，导致回复失败。
- 目前仅测试过 gpt-3.5-turbo 模型的支持。
