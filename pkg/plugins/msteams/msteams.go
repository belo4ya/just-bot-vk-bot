package msteams

import (
	"errors"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	vk "github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	"github.com/belo4ya/just-bot-vk-bot/pkg/fuapi"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

type Plugin struct {
	subscribers []*Subscriber
}

func (p *Plugin) Apply(b *vk.Bot) {
	if err := b.DB.AutoMigrate(&Task{}); err != nil {
		log.Fatalln(err)
	}

	s := NewSubscriber(b)
	p.subscribers = append(p.subscribers, s)

	//_, err := b.AddCronTask("@every 30s", s.CronTask())
	_, err := b.AddCronTask("55 00 * * *", s.CronTask()) // "55 03 * * *" (UTC+3)
	if err != nil {
		log.Fatalln(err)
	}
}

type Subscriber struct {
	ChatID    int
	GroupName string
	Delay     time.Duration
	bot       *vk.Bot
}

func NewSubscriber(b *vk.Bot) *Subscriber {
	return &Subscriber{
		ChatID:    2000000001,
		GroupName: "ПИ19-3",
		Delay:     -5 * time.Minute,
		bot:       b,
	}
}

func (s *Subscriber) CronTask() vk.CronTask {
	return func() {
		group, err := fuapi.GetGroup(s.GroupName)
		if err != nil {
			log.Fatalln(err)
		}

		schedule, err := fuapi.GetGroupSchedule(group.ID, time.Now().UTC(), time.Now().UTC())
		if err != nil {
			log.Fatalln(err)
		}

		for _, item := range *schedule {
			if item.URL1 == "" {
				continue
			}

			startAt := s.startAt(item.BeginLesson)
			sendAt := s.sendAt(startAt)
			dur := s.duration(sendAt)
			if dur < 0 {
				continue
			}

			task := &Task{Status: PENDING, SendAt: sendAt, Message: Message(item), StartAt: startAt}
			if err := s.bot.DB.Create(task).Error; err != nil {
				log.Fatalln(err)
			}

			time.AfterFunc(dur, s.Task(task.ID))
			log.Printf("Created Task{ID: %d, SendAt: %s}, dur = %s", task.ID, task.SendAt.Format(time.RFC3339), dur)
		}

		p := params.NewMessagesSendBuilder()
		p.Message("Ok").PeerID(s.ChatID).RandomID(vk.RandomID())

		if _, err := s.bot.VK.MessagesSend(p.Params); err != nil {
			log.Fatalln(err)
		}
	}
}

func (s *Subscriber) Task(taskID uint) func() {
	return func() {
		var task Task
		err := s.bot.DB.First(&task, taskID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return
			}
			log.Fatalln(err)
		}

		if task.Status != PENDING {
			return
		}

		p := params.NewMessagesSendBuilder()
		p.Message(fmt.Sprintf(task.Message)).PeerID(s.ChatID).RandomID(vk.RandomID())
		if _, err := s.bot.VK.MessagesSend(p.Params); err != nil {
			log.Fatalln(err)
		}
		s.bot.DB.Model(task).Update("status", SUCCESS)
	}
}

func (s *Subscriber) startAt(beginLesson string) time.Time {
	var start [2]int
	for i, s := range strings.Split(beginLesson, ":") {
		start[i], _ = strconv.Atoi(s)
	}
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), start[0], start[1], 0, 0, time.Local)
}

func (s *Subscriber) sendAt(startAt time.Time) time.Time {
	return startAt.Add(s.Delay)
}

func (s *Subscriber) duration(sendAt time.Time) time.Duration {
	return sendAt.Sub(time.Now().UTC())
	//return time.Duration(rand.Intn(20-3)+3) * time.Second
}
