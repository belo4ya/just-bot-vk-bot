package main

import (
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
	vkAdapter := b.Adapter.(*joeVk.BotAdapter)

	c := cron.New()
	entryID, err := c.AddFunc("@every 15s", func() {
		if err := vkAdapter.Send(time.Now().Local().Format("15:04 02.01.2006"), "510253487"); err != nil {
			b.Logger.Fatal(err.Error())
		}
	})
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
	log.Printf("EntryID: %d", entryID)
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
