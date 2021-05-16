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
	})

	b.Start()
}

func SendCutePhoto(animal int, to *tb.Chat, b *tb.Bot) {
	var getCuteAnimal func() *tb.Photo

	switch animal {
	case DOG:
		getCuteAnimal = getRandomDog
	case CAT:
		getCuteAnimal = getRandomCat
	case GUINEA_PIG:
		getCuteAnimal = getRandomGuineaPig

	}
	animalPhoto, success := tryHard(getCuteAnimal, 10)
	if success {
		animalPhoto.Send(b, to, &tb.SendOptions{})
	}
}
