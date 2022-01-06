package plugins

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	vk "github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	"github.com/belo4ya/just-bot-vk-bot/pkg/fuapi"
	"gorm.io/gorm"
	"log"
	"math/rand"
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

func (s *Subscriber) Task() vk.CronTask {
	return func(b *vk.Bot) func() {
		return func() {
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

			for _, item := range *schedule {
				if item.URL1 != "" {
					dur := time.Duration(rand.Intn(20-3)+3) * time.Second
					msg := fmt.Sprintf("⏱%s - %s⏱\n", item.BeginLesson, item.EndLesson) +
						item.Discipline + "\n" +
						item.KindOfWork + "\n" +
						"Кто: " + item.Lecturer + "\n" +
						"Где: " + item.URL1 + "\n" +
						choiceNote(notes)

					time.AfterFunc(dur, func() {
						p := params.NewMessagesSendBuilder()
						p.Message(fmt.Sprintf(msg)).PeerID(s.ChatID).RandomID(vk.RandomID())
						if _, err := b.VK.MessagesSend(p.Params); err != nil {
							log.Fatalln(err)
						}
					})

					if err := s.repo.SaveTask(&Task{SendAt: time.Now().Add(dur), Message: msg}); err != nil {
						log.Fatalln(err)
					}
				}
			}

			p := params.NewMessagesSendBuilder()
			p.Message("Ok").PeerID(s.ChatID).RandomID(vk.RandomID())

			if _, err := b.VK.MessagesSend(p.Params); err != nil {
				log.Fatalln(err)
			}
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

func TeamsInit(b *vk.Bot) {
	if err := b.DB.AutoMigrate(&Task{}); err != nil {
		log.Fatalln(err)
	}

	_, err := b.AddCronTask("55 03 * * *", NewSubscriber(NewTaskRepo(b.DB)).Task())
	if err != nil {
		log.Fatalln(err)
	}
}
