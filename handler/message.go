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
	dispatcher := openwechat.NewMessageMatchDispatcher()

	dispatcher.OnText(func(msgCtx *openwechat.MessageContext) {
		msg := msgCtx.Message

		if len(msg.Content) <= 4 {
			return
		}
		if strings.Index(msg.Content, " ") == 0 {

		}
		if msg.Content[:4] == ":bot" {
			query := msg.Content[4:]

			response, err := model.ChatCompletion(ctx, model.MakeMessage(query))
			if err != nil {
				log.Println("ChatCompletion error: " + err.Error())
				return
			}
			log.Println("回复消息: " + response)

			receiver, err := msg.Receiver()
			if err != nil {
				log.Println("Receiver error: " + err.Error())
				return
			}
			log.Printf("receiver: %v, isSendBySelf: %v", receiver, msg.IsSendBySelf())
			sender, err := msg.Sender()
			if err != nil {
				log.Println("Sender error: " + err.Error())
				return
			}
			log.Printf("sender: %v, isSendBySelf: %v", sender, msg.IsSendBySelf())

			if receiver != nil && receiver.UserName == "filehelper" {
				log.Println("Reply to filehelper")
				fh := self.FileHelper()
				_, err := self.SendTextToFriend(fh, response)
				if err != nil {
					log.Println("SendTextToFriend error: " + err.Error())
					return
				}
			} else if receiver != nil && msg.IsSendBySelf() {
				log.Println("Is send by self, isGroup:", receiver.IsGroup(), ", isSelf:", receiver.IsSelf(), ", isFriend:", receiver.IsFriend())
				if receiver.IsGroup() {
					group, _ := receiver.AsGroup()
					log.Println("Reply to group")
					_, err := self.SendTextToGroup(group, response)
					if err != nil {
						log.Println("SendTextToGroup error: " + err.Error())
					}
				} else {
					user, _ := receiver.AsFriend()
					log.Println("Reply to user")
					_, err := self.SendTextToFriend(user, response)
					if err != nil {
						log.Println("SendTextToFriend error: " + err.Error())
					}
				}
			} else {
				log.Println("Fallback reply to message")
				_, err := msgCtx.ReplyText(response)
				if err != nil {
					log.Println("ReplyText error: " + err.Error())
				}
			}
		}
	})

	return dispatcher.AsMessageHandler()
}
