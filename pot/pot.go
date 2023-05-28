package pot

type Pot struct {
	Level        int64 `json:"level"`
	Wager        int64 `json:"wager"`
	Total        int64 `json:"total"`
	Contributors []int `json:"contributors"`
}

func (p *Pot) ContributorExists(idx int) bool {
	for _, cIdx := range p.Contributors {
		if cIdx == idx {
			return true
		}
	}

	return false
}
