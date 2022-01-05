package main

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/go-joe/joe"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"time"

	joeVk "github.com/tdakkota/joe-vk-adapter"
)

func registerHandlers(b *joe.Bot) {
	b.Respond("ping", func(msg joe.Message) error {
		msg.Respond("pong")
		return nil
	})
}

func registerCron(b *joe.Bot) {
	task1 := func() {
		vk := api.NewVK(os.Getenv("TOKEN"))
		now := time.Now().Local().Format("15:04 02.01.2006")
		msg := fmt.Sprintf("%s: task 1", now)
		params := api.Params{
			"peer_id":   510253487,
			"message":   msg,
			"random_id": 0,
		}
		if _, err := vk.MessagesSend(params); err != nil {
			b.Logger.Fatal(err.Error())
		}
	}

	task2 := func() {
		vkAdapter := b.Adapter.(*joeVk.BotAdapter)
		now := time.Now().Local().Format("15:04 02.01.2006")
		msg := fmt.Sprintf("%s: task 2", now)
		if err := vkAdapter.Send(msg, "510253487"); err != nil {
			b.Logger.Fatal(err.Error())
		}
	}

	c := cron.New()
	c.AddFunc("@every 5s", task1)
	c.AddFunc("@every 15s", task2)
	c.Start()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	log.Printf("TOKEN: %s", token)

	b := joe.New("just-bot", joeVk.Adapter(token))

	registerCron(b)
	registerHandlers(b)

	if err := b.Run(); err != nil {
		b.Logger.Fatal(err.Error())
	}
}
