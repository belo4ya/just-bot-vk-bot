package plugins

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	"github.com/belo4ya/just-bot-vk-bot/pkg/fu-api"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	chatID    int    = 2000000001
	groupName string = "ПИ19-3"
)

var (
	notes = []string{
		"Ребятки, не опаздывайте",
		"Коллеги, убедительная просьба не опаздывать",
	}
)

type Subscriber struct {
	ChatID    int
	GroupName string
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		ChatID:    chatID,
		GroupName: groupName,
	}
}

func (s *Subscriber) Handler() bot.Handler {
	return func(vk *api.VK, obj events.MessageNewObject) {
		group, err := fuapi.GetGroup(s.GroupName)
		if err != nil {
			log.Fatalln(err)
		}

		schedule, err := fuapi.GetGroupSchedule(
			group.ID,
			time.Date(2021, 12, 27, 0, 0, 0, 0, time.Local),
			time.Date(2021, 12, 27, 0, 0, 0, 0, time.Local),
		)
		if err != nil {
			log.Fatalln(err)
		}

		var b strings.Builder
		for _, item := range *schedule {
			if item.URL1 != "" {
				b.WriteString(fmt.Sprintf("⏱%s - %s⏱\n", item.BeginLesson, item.EndLesson))
				b.WriteString(item.Discipline + "\n")
				b.WriteString(item.KindOfWork + "\n")
				b.WriteString("Кто: " + item.Lecturer + "\n")
				b.WriteString("Где: " + item.URL1 + "\n")
				b.WriteString(choiceNote(notes) + "\n\n")
			}
		}

		p := params.NewMessagesSendBuilder()
		p.Message(b.String()).PeerID(s.ChatID).RandomID(bot.RandomID())

		if _, err := vk.MessagesSend(p.Params); err != nil {
			log.Fatalln(err)
		}
	}
}

func choiceNote(notes []string) string {
	return notes[rand.Intn(len(notes))]
}
