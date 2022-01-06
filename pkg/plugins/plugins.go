package plugins

import (
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	"log"
)

func PingHandler(b *bot.Bot, obj events.MessageNewObject) {
	p := params.NewMessagesSendBuilder()
	p.Message("Everything all right ☕").PeerID(obj.Message.PeerID).RandomID(bot.RandomID())
	if _, err := b.VK.MessagesSend(p.Params); err != nil {
		log.Fatalln(err)
	}
}

func HelloHandler(b *bot.Bot, obj events.MessageNewObject) {
	p := params.NewMessagesSendBuilder()
	p.Message("👋").PeerID(obj.Message.PeerID).RandomID(bot.RandomID())
	if _, err := b.VK.MessagesSend(p.Params); err != nil {
		log.Fatalln(err)
	}
}
