package combination

import "sort"

type Element struct {
	Combination Combination
	Rank        int
	Count       int
}

func GetElementsByRank(cards []*Card) []*Element {

	result := make(map[int]*Element)
	elements := make([]*Element, 0)

	for _, c := range cards {
		_, ok := result[c.Rank]
		if ok {
			ele := result[c.Rank]
			ele.Count++

			if ele.Count == 2 {
				ele.Combination = CombinationPair
			} else if ele.Count == 3 {
				ele.Combination = CombinationThreeOfAKind
			} else if ele.Count == 4 {
				ele.Combination = CombinationFourOfAKind
			}

			continue
		}

		ele := &Element{
			Rank:        c.Rank,
			Combination: CombinationHighCard,
			Count:       1,
		}

		result[c.Rank] = ele
		elements = append(elements, ele)
	}

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Count > elements[j].Count
	})

	return elements
}
