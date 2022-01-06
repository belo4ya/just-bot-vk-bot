package plugins

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	vkBot "github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	"github.com/belo4ya/just-bot-vk-bot/pkg/fu-api"
	"gorm.io/gorm"
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

func choiceNote(notes []string) string {
	return notes[rand.Intn(len(notes))]
}

type Subscriber struct {
	ChatID    int
	GroupName string
	repo      *TaskRepo
}

func NewSubscriber(repo *TaskRepo) *Subscriber {
	return &Subscriber{
		ChatID:    chatID,
		GroupName: groupName,
		repo:      repo,
	}
}

func (s *Subscriber) Handler() vkBot.Handler {
	return func(bot *vkBot.Bot, obj events.MessageNewObject) {
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
		p.Message(b.String()).PeerID(s.ChatID).RandomID(vkBot.RandomID())

		if err := s.repo.SaveTask(&Task{SendAt: time.Now(), Message: b.String()}); err != nil {
			log.Fatalln(err)
		}

		if _, err := bot.VK.MessagesSend(p.Params); err != nil {
			log.Fatalln(err)
		}
	}
}

type Task struct {
	gorm.Model
	SendAt  time.Time
	Message string
}

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r TaskRepo) SaveTask(t *Task) error {
	return r.db.Create(t).Error
}

func TeamsInit(b *vkBot.Bot) {
	if err := b.DB.AutoMigrate(&Task{}); err != nil {
		log.Fatalln(err)
	}
	taskRepo := NewTaskRepo(b.DB)
	b.AddHandler("schedule", NewSubscriber(taskRepo).Handler())
}
