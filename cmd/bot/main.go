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
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TOKEN")

	b := bot.NewBot(token)

	b.AddHandler("ping", plugins.PingHandler)
	b.AddHandler("hello", plugins.HelloHandler)
	if _, err := b.AddCronTask("@every 5s", plugins.CronTask); err != nil {
		log.Fatal(err)
	}

	log.Println("Start Long Poll")
	if err := b.Run(); err != nil {
		log.Fatal(err)
	}
}
