package main

import (
	"fmt"
	"math/rand"
	"time"
)

type CardSuit int32

const (
	CardSuitSpade CardSuit = iota
	CardSuitHeart
	CardSuitDiamond
	CardSuitClub
)

var CardSuits = []string{
	"S",
	"H",
	"D",
	"C",
}

var CardPoints = []string{
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"T",
	"J",
	"Q",
	"K",
	"A",
}

func NewStandardDeckCards() []string {

	cards := make([]string, 0, 52)

	for _, suit := range CardSuits {
		for i := 0; i < 13; i++ {
			cards = append(cards, fmt.Sprintf("%s%s", suit, CardPoints[i]))
		}
	}

	return cards
}

func NewShortDeckCards() []string {

	cards := make([]string, 0, 36)

	for _, suit := range CardSuits {

		// Take off 2, 3, 4 and 5
		for i := 4; i < 13; i++ {
			cards = append(cards, fmt.Sprintf("%s%s", suit, CardPoints[i]))
		}
	}

	return cards
}

func ShuffleCards(cards []string) []string {

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}
