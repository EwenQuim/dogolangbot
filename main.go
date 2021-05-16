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
	count  int
	emoji  string
	winner bool
}

type Dogobot struct {
	animals     map[string]*Animal
	total_calls int
}

var dogobot Dogobot

func init() {
	dogobot = Dogobot{
		animals: map[string]*Animal{
			"woof":  {emoji: "🐶"},
			"meow":  {emoji: "🐱"},
			"pouic": {emoji: "🐷🇮🇳"},
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

	b.Handle("/woof", func(m *tb.Message) {
		destinataire := m.Chat
		go SendCutePhoto(DOG, destinataire, b)
	})

	b.Handle("/meow", func(m *tb.Message) {
		destinataire := m.Chat
		go SendCutePhoto(CAT, destinataire, b)
	})

	b.Handle("/pouic", func(m *tb.Message) {
		destinataire := m.Chat
		go SendCutePhoto(GUINEA_PIG, destinataire, b)
	})

	b.Handle("/winner", func(m *tb.Message) {

		response := dogobot.updateScores()
		b.Send(m.Chat, response)
	})

	b.Start()
}
