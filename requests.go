package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (dogobot Dogobot) SendCutePhoto(message string, to *tb.Chat, b *tb.Bot) error {
	slog.Info("received", "message", message)

	messageSplit := strings.Fields(message)

	if len(messageSplit) == 0 {
		return nil
	}

	animal := strings.TrimLeft(messageSplit[0], "/")
	if _, exists := dogobot.animals[animal]; !exists {
		return errors.New("unknown command " + animal)
	}

	animalPhoto := tryHard(dogobot.animals[animal].function, 10)
	if animalPhoto == nil {
		return errors.New("image not found after 10 tries")
	}

	_, err := animalPhoto.Send(b, to, &tb.SendOptions{})
	if err != nil {
		return err
	}

	go func() {
		err = dogobot.saveToDatabase(animal)
		if err != nil {
			slog.Error("error saving to database", "error", err)
		}
	}()

	return nil
}

var netClient = &http.Client{
	Timeout: time.Second * 3,
}

func getRandomDog() *tb.Photo {
	// http request to the API
	resp, err := netClient.Get("https://random.dog/woof")
	if err != nil {
		if os.IsTimeout(err) {
			log.Fatalf("timeout getting a dog image")
		}
		log.Fatalf("error getting a dog image: %v", err)
	}
	defer resp.Body.Close()

	// decode photo url sent
	var result []byte
	result, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error decoding a dog image: %v", err)
	}
	// result := "yo.jpg"
	photoUrl := "https://random.dog/" + strings.ToLower(string(result))
	fmt.Println("found dog image:", photoUrl)

	ext := filepath.Ext(photoUrl)
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		return &tb.Photo{File: tb.FromURL(photoUrl)}
	}
	return nil
}

func getRandomCat() *tb.Photo {
	// http request to the API
	resp, err := netClient.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// decode photo url sent
	var result []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}
	photoUrl, ok := result[0]["url"].(string)
	if !ok {
		return nil
	}

	ext := filepath.Ext(photoUrl)
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		return &tb.Photo{File: tb.FromURL(photoUrl)}
	}
	return nil
}

// tryHard sends n requests and returns the first satisfying result
func tryHard(f func() *tb.Photo, maxTries int) *tb.Photo {

	firstPhoto := make(chan *tb.Photo, maxTries)
	done := make(chan bool, maxTries)
	gogogo := make(chan int, maxTries)

	for i := range maxTries {
		go func(i int, done <-chan bool, gogogo chan<- int) {
			defer func() {
				if r := recover(); r != nil && os.Getenv("ENV") == "dev" {
					fmt.Println("recovered:", r)
				}
			}()

			time.Sleep(time.Duration(i) * time.Duration(i) * 10 * time.Millisecond)

			gogogo <- i

		}(i, done, gogogo)
	}

	for {
		select {
		case photo := <-firstPhoto:
			return photo
		case <-gogogo:
			go func(firstPhoto chan<- *tb.Photo) {
				defer func() {
					if r := recover(); r != nil && os.Getenv("ENV") == "dev" {
						fmt.Println("recovered:", r)
					}
				}()
				photo := f()
				if photo != nil {
					firstPhoto <- photo
					close(gogogo)
				}
			}(firstPhoto)
		}

	}

}
