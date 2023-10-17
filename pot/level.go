package pot

type Level struct {
	Level        int64 `json:"level"`
	Wager        int64 `json:"wager"`
	Total        int64 `json:"total"`
	Contributors []int `json:"contributors"`
}

func (l *Level) ContributorExists(idx int) bool {
	for _, cIdx := range l.Contributors {
		if cIdx == idx {
			return true
		}
	}

	return false
}
