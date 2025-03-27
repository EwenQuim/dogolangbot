package main

import (
	"testing"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func TestTryHardDog(t *testing.T) {
	photo := tryHard(getRandomDog, 5)
	t.Log("photo obtained", photo)
	if photo == nil {
		t.Fail()
	}
}

func TestTryHardCat(t *testing.T) {
	photo := tryHard(getRandomCat, 5)
	t.Log("photo obtained", photo)
	if photo == nil {
		t.Fail()
	}
}

func TestTryHard(t *testing.T) {

	t.Run("succeeds because under 2500ms = 5*5*10ms", func(t *testing.T) {
		photo := tryHard(func() *tb.Photo {
			time.Sleep(100 * time.Millisecond)
			return &tb.Photo{File: tb.File{FileID: "123"}}
		}, 5)
		t.Log("photo not obtained ", photo)
		if photo == nil {
			t.Fail()
		}
	})

	t.Run("fail because over 2500ms = 5*5*10ms", func(t *testing.T) {
		photo := tryHard(func() *tb.Photo {
			time.Sleep(3000 * time.Millisecond)
			return &tb.Photo{File: tb.File{FileID: "123"}}
		}, 5)
		t.Log("photo obtained while shouldn't", photo)
		if photo != nil {
			t.Fail()
		}
	})

	t.Run("always fail should block nor leak", func(t *testing.T) {
		photo := tryHard(func() *tb.Photo { return nil }, 5)
		t.Log("photo obtained while shouldn't", photo)
		if photo != nil {
			t.Fail()
		}
	})
}
