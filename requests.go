package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		animal = messageSplit[0][1:]
		if _, exists := dogobot.animals[animal]; !exists {
			fmt.Println("unknown command, aborting")
			return nil
		}
	} else {
		return nil
	}

	var success bool
	var animalPhoto *tb.Photo
	if dogobot.animals[animal].subreddit == "" {
		animalPhoto, success = tryHard(dogobot.animals[animal].function, 10)
	} else {
		animalPhoto, success = tryHard(func() *tb.Photo { return getFromReddit(dogobot.animals[animal].subreddit) }, 10)
	}

	if success {
		_, err := animalPhoto.Send(b, to, &tb.SendOptions{})
		if err == nil {
			saveToDatabase(animal)
			return nil
		}
	}
	fmt.Println("image not found after 10 tries")
	return nil
}

func getRandomDog() *tb.Photo {
	// http request to the API
	resp, err := http.Get("https://random.dog/woof.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// decode photo url sent
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	photoUrl := strings.ToLower(result["url"])

	for _, fileType := range []string{"jpg", "peg", "png"} {
		if photoUrl[len(photoUrl)-3:] == fileType {
			return &tb.Photo{File: tb.FromURL(photoUrl)}
		}
	}
	return nil
}

func getRandomCat() *tb.Photo {
	// http request to the API
	resp, err := http.Get("http://aws.random.cat/meow")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// decode photo url sent
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	photoUrl := result["file"]

	for _, fileType := range []string{"jpg", "jpeg", "png"} {
		if photoUrl[len(photoUrl)-4:] == fileType {
			return &tb.Photo{File: tb.FromURL(photoUrl)}
		}
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

// func getRandomGuineaPig() *tb.Photo {
// 	return getFromReddit("guineapigs")
// }

func getFromReddit(subreddit string) *tb.Photo {

	client := &http.Client{}
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
	jason, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err)
	}

	// decode photo url sent
	var result []Thing
	err = json.Unmarshal(jason, &result)
	if err != nil {
		// do anything
		println(err)
	}

	photoUrl := result[0].Data.Children[0].Data.UrlOverriddenByDest
	fmt.Printf("photo link: %v\n", photoUrl)

	if photoUrl == "" || photoUrl[:17] != "https://i.redd.it" {
		return nil
	}
	return &tb.Photo{File: tb.FromURL(photoUrl)}

}

func tryHard(f func() *tb.Photo, maxTries int) (*tb.Photo, bool) {

	for tries := 0; tries < maxTries; tries++ {

		photo := f()
		if photo != nil {
			return photo, true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil, false
}

// func tryHard(tries int, f func() *tb.Photo, photo chan *tb.Photo) {

// 	for i := 0; i < tries; i++ {

// 		go func() {
// 			phoTry := f()
// 			if photo != nil {
// 				photo <- phoTry
// 			}
// 		}()

// 		time.Sleep(100 * time.Millisecond)
// 	}
// }
