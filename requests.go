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

	if animalPhoto != nil {
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
		ext := strings.Split(photoUrl, ".")
		if ext[len(ext)-1] == fileType {
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
		ext := strings.Split(photoUrl, ".")
		if ext[len(ext)-1] == fileType {
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
	for try := 0; try < maxTries; try++ {
		go func(try int, firstPhoto chan<- *tb.Photo) {
			time.Sleep(time.Duration(try) * time.Duration(try) * 10 * time.Millisecond)
			photo := f()
			if photo != nil {
				firstPhoto <- photo
			}
		}(try, firstPhoto)
	}

	return <-firstPhoto
}
