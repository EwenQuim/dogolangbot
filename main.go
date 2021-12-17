package main

import (
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	tb "gopkg.in/tucnak/telebot.v2"
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

func main() {

	dogobot := Dogobot{
		animals: map[string]*Animal{
			"woof":  {emoji: "🐶", function: getRandomDog},
			"meow":  {emoji: "🐱", function: getRandomCat},
			"pouic": {emoji: "🐹", subreddit: "guineapigs"},
			"awww":  {emoji: "🥰", subreddit: "awww"},
			"earth": {emoji: "🌍", subreddit: "earthPorn"},
		},
		total_calls: 0,
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	// Handle any command not already handled that begins by `/`
	b.Handle(tb.OnText, func(m *tb.Message) {
		if m.Text[0] == '/' {
			destinataire := m.Chat
			go dogobot.SendCutePhoto(m.Text, destinataire, b)
		}
	})

	b.Handle("/winner", func(m *tb.Message) {
		go b.Send(m.Chat, dogobot.getScores())
	})

	b.Start()
}
