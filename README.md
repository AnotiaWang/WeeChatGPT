# WeeChatGPT

将 ChatGPT 集成到微信个人号。基于 [openwechat](https://github.com/eatmoreapple/openwechat)，支持网页版和桌面版微信。

## 运行方法

1. 前往 [Releases](https://github.com/AnotiaWang/WeeChatGPT/releases/latest) ，根据你使用的平台，下载最新版本的 WeeChatGPT。
2. 先运行一次 WeeChatGPT，程序会生成配置文件 `config.yml`，然后根据文件中的提示，填写配置。
3. 再次运行 WeeChatGPT，会显示微信的登录二维码链接，在浏览器中打开它，扫码登录微信即可。

![](https://i.328888.xyz/2023/03/13/v602P.png)

### 说明

- 程序支持热登录，如果两次登录之间的间隔较短，可以自动登录。
- 建议将群组添加到通讯录，否则机器人可能无法获取到群组信息，导致回复失败。
- 目前仅测试过 gpt-3.5-turbo 模型的支持。
- 目前不支持多轮对话。

如果无法访问 OpenAI 的 API，可以考虑我自建的 Nginx 反代 https://proxy.api.ataw.top/openai/v1/chat/completions 。没有暗箱操作，请放心使用。

<details>
<summary>Nginx 反代配置</summary>

```text
server {
    listen       80;
    listen       443 ssl http2;
    server_name  proxy.api.ataw.top;

    ssl_certificate /etc/nginx/conf.d/proxy.api.ataw.top_bundle.crt; 
    ssl_certificate_key /etc/nginx/conf.d/proxy.api.ataw.top.key; 
    ssl_session_timeout 5m;
    ssl_protocols TLSv1.2 TLSv1.3; 
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE; 
    ssl_prefer_server_ciphers on;

    location /openai/ {
        proxy_pass https://api.openai.com/;
        proxy_set_header Host api.openai.com;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```
</details>
