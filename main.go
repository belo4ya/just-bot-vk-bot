package main

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"time"
)

type VkBot struct {
	vk *api.VK
	lp *longpoll.LongPoll
	c  *cron.Cron
}

func NewVkBot(token string) *VkBot {
	log.Printf("TOKEN: %s", token)
	vk := api.NewVK(token)

	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	c := cron.New()

	return &VkBot{
		vk: vk,
		lp: lp,
		c:  c,
	}
}

func (b *VkBot) AddHandler(pattern string, handler func(*VkBot, events.MessageNewObject)) {
	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		if obj.Message.Text == pattern {
			handler(b, obj)
		}
	})
}

func (b *VkBot) AddCronTask(spec string, task func(*VkBot) func()) (cron.EntryID, error) {
	cmd := task(b)
	return b.c.AddFunc(spec, cmd)
}

func (b *VkBot) Run() error {
	b.c.Start()
	return b.lp.Run()
}

func PingHandler(b *VkBot, obj events.MessageNewObject) {
	p := params.NewMessagesSendBuilder()
	p.Message("pong")
	p.RandomID(0)
	p.PeerID(obj.Message.PeerID)

	if _, err := b.vk.MessagesSend(p.Params); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler(b *VkBot, obj events.MessageNewObject) {
	p := params.NewMessagesSendBuilder()
	p.Message("Hi !")
	p.RandomID(0)
	p.PeerID(obj.Message.PeerID)

	if _, err := b.vk.MessagesSend(p.Params); err != nil {
		log.Fatal(err)
	}
}

func CronTask(b *VkBot) func() {
	return func() {
		now := time.Now().Local().Format("15:04 02.01.2006")
		msg := fmt.Sprintf("%s: task 1", now)

		p := params.NewMessagesSendBuilder()
		p.PeerID(510253487)
		p.Message(msg)
		p.RandomID(0)

		if _, err := b.vk.MessagesSend(p.Params); err != nil {
			log.Fatal(err)
		}
	}
}

//func registerCron(b *joe.Bot) {
//	task1 := func() {
//		vk := api.NewVK(os.Getenv("TOKEN"))
//
//		p := params.NewMessagesSendBuilder()
//		conversations, err := vk.MessagesGetConversations(p.Params)
//		if err != nil {
//			b.Logger.Fatal(err.Error())
//		}
//		b.Logger.Info(fmt.Sprintf("Conversations: %+v", conversations))
//
//		now := time.Now().Local().Format("15:04 02.01.2006")
//		msg := fmt.Sprintf("%s: task 1", now)
//		p = params.NewMessagesSendBuilder()
//		p.PeerID(510253487)
//		p.Message(msg)
//		p.RandomID(0)
//		if _, err := vk.MessagesSend(p.Params); err != nil {
//			b.Logger.Fatal(err.Error())
//		}
//	}
//
//	c := cron.New()
//	c.AddFunc("@every 5s", task1)
//	c.Start()
//}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TOKEN")

	b := NewVkBot(token)

	b.AddHandler("ping", PingHandler)
	b.AddHandler("hello", HelloHandler)
	if _, err := b.AddCronTask("@every 5s", CronTask); err != nil {
		log.Fatal(err)
	}

	log.Println("Start Long Poll")
	if err := b.Run(); err != nil {
		log.Fatal(err)
	}
}
