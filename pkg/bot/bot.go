package bot

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/robfig/cron/v3"
	"log"
)

type (
	Handler  func(*Bot, events.MessageNewObject)
	CronTask func(*Bot) func()
)

type Bot struct {
	Vk *api.VK
	lp *longpoll.LongPoll
	c  *cron.Cron
}

func NewBot(token string) *Bot {
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

	return &Bot{
		Vk: vk,
		lp: lp,
		c:  c,
	}
}

func (b *Bot) AddHandler(pattern string, handler Handler) {
	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		if obj.Message.Text == pattern {
			handler(b, obj)
		}
	})
}

func (b *Bot) AddCronTask(spec string, task CronTask) (cron.EntryID, error) {
	cmd := task(b)
	return b.c.AddFunc(spec, cmd)
}

func (b *Bot) Run() error {
	b.c.Start()
	return b.lp.Run()
}
