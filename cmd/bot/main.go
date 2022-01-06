package main

import (
	"github.com/belo4ya/just-bot-vk-bot/pkg/bot"
	"github.com/belo4ya/just-bot-vk-bot/pkg/plugins"
	"github.com/belo4ya/just-bot-vk-bot/pkg/plugins/msteams"
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

	b.AddPlugin(&msteams.MsTeams{})

	log.Println("Start Long Poll")
	if err := b.Run(); err != nil {
		log.Fatalln(err)
	}
}
