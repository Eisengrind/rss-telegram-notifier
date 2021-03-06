package main

import (
	"bytes"
	"log"
	"time"

	"github.com/playnet-public/flagenv"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	tgAPIKey   = flagenv.String("tg-api-key", "", "This is your telegram bot-api-token")
	chatID     = flagenv.Int("chat-id", 0, "The chat to use.")
	rssFeedURL = flagenv.String("rss-feed-url", "", "The RSS-feed to use for updates")
	timeout    = flagenv.Int("timeout", 10, "The timeout within a rss feed is gathered (in seconds)")
)

type message struct {
	Title       string
	Description string
	Date        string
	Link        string
	Author      string
}

func main() {
	log.Println("starting rss feed puller")
	log.Println("parsing arguments")
	flagenv.Parse()
	log.Println("arguments parsed")

	if *chatID == 0 {
		log.Fatal("you have to enter an chat id")
	}

	if *rssFeedURL == "" {
		log.Fatal("you have to enter a RSS feed url")
	}

	log.Println("creating new telegram bot")
	bot, err := tgbotapi.NewBotAPI(*tgAPIKey)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("telegram bot created")

	log.Println("creating new rss feed puller")
	updator := NewRSSFeed(*rssFeedURL, time.Duration(*timeout)*time.Second)
	updates := updator.GetUpdateChan()

	log.Println("rss feed puller created\nwaiting for new items")
	for update := range updates {
		var b bytes.Buffer

		if err := parsedTemplateMessage.Execute(&b, &message{
			Title:       update.Title,
			Description: update.Description,
			Date:        update.PublishedParsed.Format("02.01.2006 @ 15:04"),
			Link:        update.Link,
			Author:      update.Author.Name,
		}); err != nil {
			log.Fatal(err.Error())
		}

		msg := tgbotapi.NewMessage(
			int64(*chatID),
			b.String(),
		)
		msg.ParseMode = "html"

		if _, err = bot.Send(msg); err != nil {
			log.Fatal(err.Error())
		}
	}
}
