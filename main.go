package main

import (
	"database/sql"
	"expvar"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	_ "github.com/joho/godotenv/autoload"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Animal struct {
	count    int
	emoji    string
	function func() *tb.Photo // leave subreddit field empty is function is set to something
}

func goroutines() any {
	return runtime.NumGoroutine()
}

type Dogobot struct {
	animals     map[string]*Animal
	total_calls int
	db          *sql.DB
}

func main() {

	db := createDb()

	dogobot := Dogobot{
		animals: map[string]*Animal{
			"woof": {emoji: "🐶", function: getRandomDog},
			"meow": {emoji: "🐱", function: getRandomCat},
		},
		total_calls: 0,
		db:          db,
	}

	fmt.Println("Starting bot")

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	for command := range dogobot.animals {
		slog.Info("registering command", "command", command)
		b.Handle("/"+command, func(m *tb.Message) {
			err := dogobot.SendCutePhoto(command, m.Chat, b)
			if err != nil {
				slog.Error("error sending photo", "command", command, "error", err)
			}
		})
	}

	slog.Info("registering command", "command", "winner")
	b.Handle("/winner", func(m *tb.Message) {
		_, err := b.Send(m.Chat, dogobot.getScores())
		if err != nil {
			panic(err)
		}
	})

	if os.Getenv("ENV") == "dev" {
		expvar.Publish("Goroutines", expvar.Func(goroutines))
		go http.ListenAndServe(":1234", nil)
	}
	b.Start()

}
