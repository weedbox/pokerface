package pot

type Pot struct {
	Level        int64         `json:"level"`
	Wager        int64         `json:"wager"`
	Total        int64         `json:"total"`
	Contributors map[int]int64 `json:"contributors"`
	Levels       []*Level      //      `json:"levels"`
}

func (p *Pot) ContributorExists(idx int) bool {
	if _, ok := p.Contributors[idx]; ok {
		return true
	}

	return false
}
