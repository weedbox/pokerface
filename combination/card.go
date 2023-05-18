package combination

import (
	"fmt"
)

type Card struct {
	Suit string
	Rank int
}

var CardRank = map[string]int{
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"T": 10,
	"J": 11,
	"Q": 12,
	"K": 13,
	"A": 14,
}

var CardSymbol = map[int]string{
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "T",
	11: "J",
	12: "Q",
	13: "K",
	14: "A",
}

var SuitSymbol = map[int]string{
	1: "S", // Spade
	2: "H", // Heart
	3: "D", // Diamond
	4: "C", // Club
}

func GetCardState(card string) *Card {

	c := &Card{
		Suit: card[0:1],
		Rank: CardRank[card[1:2]],
	}

	return c
}

func GetCardStates(cardSymbols []string) []*Card {

	// Transform card strings to internal structure
	cards := make([]*Card, 0, len(cardSymbols))
	for _, c := range cardSymbols {
		cards = append(cards, GetCardState(c))
	}

	return cards
}

func (c *Card) ToString() string {
	return fmt.Sprintf("%s%s", c.Suit, CardSymbol[c.Rank])
}
