package main

import (
	"github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	"github.com/belo4ya/just-bot-vk-bot/pkg/plugins"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	token := os.Getenv("TOKEN")

	b := bot.NewBot(token)

	b.AddHandler("ping", plugins.PingHandler)
	b.AddHandler("hello", plugins.HelloHandler)
	b.AddHandler("schedule", plugins.NewSubscriber().Handler())

	if _, err := b.AddCronTask("@every 20s", plugins.CronTask); err != nil {
		log.Fatalln(err)
	}

	log.Println("Start Long Poll")
	if err := b.Run(); err != nil {
		log.Fatalln(err)
	}
}
