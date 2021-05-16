package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	tb "gopkg.in/tucnak/telebot.v2"
)

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
	return &tb.Photo{File: tb.FromURL(result["url"])}
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
	return &tb.Photo{File: tb.FromURL(result["file"])}
}

// thing is a Reddit type that holds all of their subtypes.
type Thing struct {
	Kind string  `json:"kind"`
	Data Listing `json:"data"`
}

type Listing struct {
	Children            []Thing `json:"children,omitempty"`
	UrlOverriddenByDest string  `json:"url_overridden_by_dest,omitempty`
}

func getRandomGuineaPig() *tb.Photo {
	client := &http.Client{}
	// http request to the Reddit API
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/guineapigs/random.json?t=all", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "android:com.coal:v1.2.0 (by /u/AmethystsStudio)")
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
	println(string(jason))

	// decode photo url sent
	var result []Thing // []map[string]map[string][]map[string]map[string]string
	err = json.Unmarshal(jason, &result)
	if err != nil {
		// do anything
		println(err)
	}

	// err = json.NewDecoder(resp.Body).Decode(&result)
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Printf("result %v\n", result)
	if len(result) >= 1 {
		return &tb.Photo{File: tb.FromURL(result[0].Data.Children[0].Data.UrlOverriddenByDest)}
	}
	return nil
}
