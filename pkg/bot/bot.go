package bot

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/robfig/cron/v3"
	"log"
	"math/rand"
	"time"
)

type (
	Handler  func(*api.VK, events.MessageNewObject)
	CronTask func(*api.VK) func()
)

type Bot struct {
	vk *api.VK
	lp *longpoll.LongPoll
	c  *cron.Cron
}

func NewBot(token string) *Bot {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	log.Printf("seed: %d", seed)

	log.Printf("token: %s", token)
	vk := api.NewVK(token)

	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatalln(err)
	}
	groupID := group[0].ID
	log.Printf("group_id: %d", groupID)

	lp, err := longpoll.NewLongPoll(vk, groupID)
	if err != nil {
		log.Fatalln(err)
	}

	c := cron.New()

	return &Bot{
		vk: vk,
		lp: lp,
		c:  c,
	}
}

func (b *Bot) AddHandler(pattern string, handler Handler) {
	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		if obj.Message.Text == pattern {
			handler(b.vk, obj)
		}
	})
}

func (b *Bot) AddCronTask(spec string, task CronTask) (cron.EntryID, error) {
	cmd := task(b.vk)
	return b.c.AddFunc(spec, cmd)
}

func (b *Bot) Run() error {
	b.c.Start()
	return b.lp.Run()
}

func RandomID() int {
	return int(rand.Int31())
}
