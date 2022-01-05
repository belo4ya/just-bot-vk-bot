package plugins

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	fuapi "github.com/belo4ya/just-bot-vk-bot/pkg/fu-api"
	"log"
	"time"
)

func PingHandler(b *bot.Bot, obj events.MessageNewObject) {
	p := params.NewMessagesSendBuilder()
	p.Message("pong")
	p.RandomID(0)
	p.PeerID(obj.Message.PeerID)

	if _, err := b.Vk.MessagesSend(p.Params); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler(b *bot.Bot, obj events.MessageNewObject) {
	p := params.NewMessagesSendBuilder()
	p.Message("Hi !")
	p.RandomID(0)
	p.PeerID(obj.Message.PeerID)

	if _, err := b.Vk.MessagesSend(p.Params); err != nil {
		log.Fatal(err)
	}
}

func ApiHandler(b *bot.Bot, obj events.MessageNewObject) {
	fuapi.GetGroup("ПИ19-3")
}

func CronTask(b *bot.Bot) func() {
	return func() {
		now := time.Now().Local().Format("15:04 02.01.2006")
		msg := fmt.Sprintf("%s: task 1", now)

		p := params.NewMessagesSendBuilder()
		p.PeerID(510253487)
		p.Message(msg)
		p.RandomID(0)

		if _, err := b.Vk.MessagesSend(p.Params); err != nil {
			log.Fatal(err)
		}
	}
}
