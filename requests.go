package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (dogobot Dogobot) SendCutePhoto(message string, to *tb.Chat, b *tb.Bot) error {
	animal := ""
	fmt.Printf("received: `%v`\n", message)
	message = strings.ReplaceAll(message, "@no_data_dog_bot", "")
	messageSplit := strings.Fields(message)

	if len(messageSplit) >= 1 {
		animal = messageSplit[0][1:] // [1:] to remove slash
		if _, exists := dogobot.animals[animal]; !exists {
			fmt.Println("unknown command, aborting")
			return nil
		}
	} else {
		return nil
	}

	var animalPhoto *tb.Photo
	if dogobot.animals[animal].subreddit == "" {
		animalPhoto = tryHard(dogobot.animals[animal].function, 10)
	} else {
		animalPhoto = tryHard(func() *tb.Photo { return getFromReddit(dogobot.animals[animal].subreddit) }, 10)
	}

	if animalPhoto == nil {
		fmt.Println("image not found after 10 tries")
		return nil
	}

	_, err := animalPhoto.Send(b, to, &tb.SendOptions{})
	if err != nil {
		return err
	}

	err = saveToDatabase(animal)
	if err != nil {
		return err
	}

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

// thing is a Reddit type that holds all of their subtypes.
type Thing struct {
	Kind string  `json:"kind"`
	Data Listing `json:"data"`
}

type Listing struct {
	Children            []Thing `json:"children,omitempty"`
	UrlOverriddenByDest string  `json:"url_overridden_by_dest,omitempty"`
}

func getFromReddit(subreddit string) *tb.Photo {

	client := &http.Client{
		Timeout: time.Second * 3,
	}
	// http request to the Reddit API
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.reddit.com/r/%v/random.json?t=all", subreddit), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "telegram:@no_data_dog_bot:v1.2.0 (by /u/AmethystsStudio)")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// read response
	jason, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err)
	}

	// decode photo url sent
	var result []Thing
	err = json.Unmarshal(jason, &result)
	if err != nil {
		fmt.Println(err)
	}
	var photoUrl string
	if len(result) > 0 && len(result[0].Data.Children) > 0 {
		photoUrl = result[0].Data.Children[0].Data.UrlOverriddenByDest
	}
	fmt.Printf("photo link: %v\n", photoUrl)

	if photoUrl == "" || len(photoUrl) < 17 || photoUrl[:17] != "https://i.redd.it" {
		return nil
	}
	return &tb.Photo{File: tb.FromURL(photoUrl)}

}

// tryHard sends n requests and returns the first satisfying result
func tryHard(f func() *tb.Photo, maxTries int) *tb.Photo {

	firstPhoto := make(chan *tb.Photo, maxTries)
	done := make(chan bool, maxTries)
	gogogo := make(chan int, maxTries)

	for i := 0; i < maxTries; i++ {
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
