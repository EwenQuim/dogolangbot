package main

import (
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	tb "gopkg.in/tucnak/telebot.v2"
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
		doggo := getRandomDog()
		_, err := doggo.Send(b, m.Chat, &tb.SendOptions{})
		if err != nil {
			log.Println(err)
		}
	})

	b.Handle("/meow", func(m *tb.Message) {
		catto := getRandomCat()
		catto.Send(b, m.Chat, &tb.SendOptions{})
	})

	b.Handle("/pouic", func(m *tb.Message) {
		malo := getRandomGuineaPig()
		if malo != nil {
			malo.Send(b, m.Chat, &tb.SendOptions{})
		}
	})

	b.Handle("/winner", func(m *tb.Message) {
	})

	b.Start()
}
