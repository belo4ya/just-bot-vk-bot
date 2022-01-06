package example

import (
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	vk "github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	"log"
)

type Plugin struct {
	bot *vk.Bot
}

func (p *Plugin) Apply(b *vk.Bot) {
	p.bot = b
	b.AddHandler("ping", p.pingHandler())
	b.AddHandler("hello", p.helloHandler())
}

func (p *Plugin) pingHandler() func(events.MessageNewObject) {
	return func(obj events.MessageNewObject) {
		pb := params.NewMessagesSendBuilder()
		pb.Message("Everything all right ☕").PeerID(obj.Message.PeerID).RandomID(vk.RandomID())
		if _, err := p.bot.VK.MessagesSend(pb.Params); err != nil {
			log.Fatalln(err)
		}
	}
}

func (p *Plugin) helloHandler() func(events.MessageNewObject) {
	return func(obj events.MessageNewObject) {
		pb := params.NewMessagesSendBuilder()
		pb.Message("👋").PeerID(obj.Message.PeerID).RandomID(vk.RandomID())
		if _, err := p.bot.VK.MessagesSend(pb.Params); err != nil {
			log.Fatalln(err)
		}
	}
}
