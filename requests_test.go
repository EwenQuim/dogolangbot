package main

import (
	"testing"

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

func TestTryHardReddit(t *testing.T) {
	photo := tryHard(func() *tb.Photo { return getFromReddit("awww") }, 5)
	t.Log("photo obtained", photo)
	if photo == nil {
		t.Fail()
	}
}
