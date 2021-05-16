package main

import (
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	DOG = iota + 1
	CAT
	GUINEA_PIG
)

type Animal struct {
	count     int
	emoji     string
	function  func() *tb.Photo // leave subreddit field empty is function is set to something
	subreddit string           // leave function field empty is subreddit is set to something
}

type Dogobot struct {
	animals     map[string]*Animal
	total_calls int
}

var dogobot Dogobot

func init() {
	dogobot = Dogobot{
		animals: map[string]*Animal{
			"woof":  {emoji: "🐶", function: getRandomDog},
			"meow":  {emoji: "🐱", function: getRandomCat},
			"pouic": {emoji: "🐷🇮🇳", subreddit: "guineapigs"},
			"awww":  {emoji: "🥰", subreddit: "awww"},
			"earth": {emoji: "🌍", subreddit: "earthPorn"},
		},
		total_calls: 0,
	}
}

func main() {

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tb.OnText, func(m *tb.Message) {
		destinataire := m.Chat
		go dogobot.SendCutePhoto(m.Text, destinataire, b)
	})

	b.Handle("/winner", func(m *tb.Message) {

		response := dogobot.updateScores()
		b.Send(m.Chat, response)
	})

	b.Start()
}
