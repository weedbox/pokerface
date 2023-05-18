package combination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCardStates(t *testing.T) {

	cardSymbols := []string{"S2", "H5", "D9", "CT", "CK"}
	cards := GetCardStates(cardSymbols)

	for i, cs := range cardSymbols {
		assert.Equal(t, cs, cards[i].ToString())
	}
}
