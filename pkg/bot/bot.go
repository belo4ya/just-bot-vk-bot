package bot

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

type (
	Handler  func(events.MessageNewObject)
	CronTask func()
	Plugin   interface {
		Apply(b *Bot)
	}
)

type Bot struct {
	VK *api.VK
	lp *longpoll.LongPoll
	c  *cron.Cron
	DB *gorm.DB
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

	c := cron.New(cron.WithLocation(time.UTC))

	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}

	return &Bot{
		VK: vk,
		lp: lp,
		c:  c,
		DB: db,
	}
}

func (b *Bot) AddPlugin(plugin Plugin) {
	plugin.Apply(b)
}

func (b *Bot) AddPlugins(plugins ...Plugin) {
	for _, p := range plugins {
		p.Apply(b)
	}
}

func (b *Bot) AddHandler(pattern string, handler Handler) {
	b.lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		if obj.Message.Text == pattern {
			handler(obj)
		}
	})
}

func (b *Bot) AddCronTask(spec string, task CronTask) (cron.EntryID, error) {
	return b.c.AddFunc(spec, task)
}

func (b *Bot) Run() error {
	b.c.Start()
	return b.lp.Run()
}

func RandomID() int {
	return int(rand.Int31())
}
