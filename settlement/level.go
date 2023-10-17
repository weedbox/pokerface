package settlement

type LevelInfo struct {
	rank Rank

	Level        int64 `json:"level"`
	Wager        int64 `json:"wager"`
	Total        int64 `json:"total"`
	Contributors []int `json:"contributors"`
}

func (li *LevelInfo) UpdateScore(playerIdx int, score int) {

	for _, c := range li.Contributors {
		if c == playerIdx {
			li.rank.AddContributor(score, playerIdx)
			break
		}
	}
}

type PotLevel struct {
	levels []*LevelInfo
}

func NewPotLevel() *PotLevel {
	return &PotLevel{
		levels: make([]*LevelInfo, 0),
	}
}

func (pl *PotLevel) AddLevel(level int64, wager int64, total int64, contributors []int) {
	pl.levels = append(pl.levels, &LevelInfo{
		Level:        level,
		Wager:        wager,
		Total:        total,
		Contributors: contributors,
	})
}
