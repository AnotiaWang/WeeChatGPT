package handler

import (
	"ChatGPT_WeChat_Bot/model"
	"context"
	"github.com/eatmoreapple/openwechat"
	"log"
	"strings"
)

func Default(ctx context.Context) openwechat.MessageHandler {
	self := ctx.Value(model.SelfKey).(*openwechat.Self)
	config := ctx.Value(model.ConfigKey).(*model.Config)
	dispatcher := openwechat.NewMessageMatchDispatcher()

	dispatcher.OnText(func(msgCtx *openwechat.MessageContext) {
		msg := msgCtx.Message

		if len(msg.Content) <= 4 {
			return
		}

		if strings.Index(msg.Content, config.OpenAI.Prefix) == 0 {
			log.Println("found message match: " + msg.Content)
			query := msg.Content[4:]
			response, err := model.ChatCompletion(ctx, model.MakeMessage(query, config.OpenAI.SystemMsg))
			if err != nil {
				log.Println("ChatCompletion error: " + err.Error())
				return
			}

			receiver, err := msg.Receiver()
			if err != nil {
				log.Println("get receiver failed: " + err.Error())
				return
			}
			log.Printf("receiver: %v, isSendBySelf: %v", receiver, msg.IsSendBySelf())
			sender, err := msg.Sender()
			if err != nil {
				log.Println("get sender failed: " + err.Error())
				return
			}
			log.Printf("sender: %v, isSendBySelf: %v", sender, msg.IsSendBySelf())

			// 和文件传输助手的消息
			if receiver != nil && receiver.UserName == "filehelper" {
				fh := self.FileHelper()
				_, err := self.SendTextToFriend(fh, response)
				if err != nil {
					log.Println("reply to filehelper failed: " + err.Error())
					return
				}
			} else if receiver != nil && msg.IsSendBySelf() {
				if receiver.IsGroup() {
					group, _ := receiver.AsGroup()
					log.Println("replying to group", group.NickName)
					_, err := self.SendTextToGroup(group, response)
					if err != nil {
						log.Println("SendTextToGroup failed: " + err.Error())
					}
				} else {
					user, _ := receiver.AsFriend()
					log.Println("replying to user", user.NickName)
					_, err := self.SendTextToFriend(user, response)
					if err != nil {
						log.Println("SendTextToFriend failed: " + err.Error())
					}
				}
			} else {
				log.Println("defaulting to replyText")
				_, err := msgCtx.ReplyText(response)
				if err != nil {
					log.Println("ReplyText failed: " + err.Error())
				}
			}
		}
	})

	return dispatcher.AsMessageHandler()
}
