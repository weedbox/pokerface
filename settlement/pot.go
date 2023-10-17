package settlement

type PotResult struct {
	rank  Rank
	level *PotLevel

	Total   int64     `json:"total"`
	Winners []*Winner `json:"winners"`
}

type Winner struct {
	Idx      int   `json:"idx"`
	Withdraw int64 `json:"withdraw"`
}

func (pr *PotResult) UpdateWinner(playerIdx int, withdraw int64) {

	for _, winner := range pr.Winners {
		if winner.Idx == playerIdx {
			winner.Withdraw += withdraw
			return
		}
	}

	w := &Winner{
		Idx:      playerIdx,
		Withdraw: withdraw,
	}

	pr.Winners = append(pr.Winners, w)

	return
}
