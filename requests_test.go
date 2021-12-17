package main

import (
	"testing"

	tb "gopkg.in/tucnak/telebot.v2"
)

func TestTryHardDog(t *testing.T) {
	photo := tryHard(getRandomDog, 10)
	t.Log("photo obtained", photo)
	if photo == nil {
		t.Fail()
	}
}

func TestTryHardReddit(t *testing.T) {
	photo := tryHard(func() *tb.Photo { return getFromReddit("awww") }, 10)
	t.Log("photo obtained", photo)
	if photo == nil {
		t.Fail()
	}
}
