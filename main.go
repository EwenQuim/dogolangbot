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
	count     int
	emoji     string
	function  func() *tb.Photo // leave subreddit field empty is function is set to something
	subreddit string           // leave function field empty is subreddit is set to something
}

func goroutines() interface{} {
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
			"woof":  {emoji: "🐶", function: getRandomDog},
			"meow":  {emoji: "🐱", function: getRandomCat},
			"pouic": {emoji: "🐹", subreddit: "guineapigs"},
			"awww":  {emoji: "🥰", subreddit: "awww"},
			"earth": {emoji: "🌍", subreddit: "earthPorn"},
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

	// Handle any command not already handled that begins by `/`
	b.Handle(tb.OnText, func(m *tb.Message) {
		if m.Text[0] == '/' {
			err := dogobot.SendCutePhoto(m.Text, m.Chat, b)
			if err != nil {
				slog.Error("error sending photo", "error", err)
			}
		}
	})

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
