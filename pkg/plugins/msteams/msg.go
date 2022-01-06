package msteams

import (
	"fmt"
	"github.com/belo4ya/just-bot-vk-bot/pkg/fuapi"
	"math/rand"
)

var (
	notes = []string{
		"Ребятки, не опаздывайте!",
		"Коллеги, убедительная просьба не опаздывать.",
	}
)

func choiceNote(notes []string) string {
	return notes[rand.Intn(len(notes))]
}

func Message(item fuapi.ScheduleItem) string {
	return fmt.Sprintf("⏱%s - %s⏱\n", item.BeginLesson, item.EndLesson) +
		item.Discipline + "\n" +
		item.KindOfWork + "\n" +
		"Кто: " + item.Lecturer + "\n" +
		"Где: " + item.URL1 + "\n" +
		choiceNote(notes)
}
