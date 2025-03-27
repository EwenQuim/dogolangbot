package main

import (
	"testing"
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
