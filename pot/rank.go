package pot

import "sort"

type RankGroup struct {
	Contributers []int
	Score        int
}

type PotRank struct {
	groups []*RankGroup
}

func NewPotRank() *PotRank {
	return &PotRank{
		groups: make([]*RankGroup, 0),
	}
}

func (pr *PotRank) AddContributer(score int, contributerIdx int) {

	for _, g := range pr.groups {

		if g.Score == score {
			g.Contributers = append(g.Contributers, contributerIdx)
			return
		}
	}

	// Not found so create a new one
	g := &RankGroup{
		Contributers: make([]int, 0),
		Score:        score,
	}

	g.Contributers = append(g.Contributers, contributerIdx)

	pr.groups = append(pr.groups, g)
}

func (pr *PotRank) Calculate() {

	// Sort by score
	sort.Slice(pr.groups, func(i, j int) bool {
		return pr.groups[i].Score > pr.groups[j].Score
	})
}

func (pr *PotRank) GetWinners() []int {

	if len(pr.groups) == 0 {
		return []int{}
	}

	return pr.groups[0].Contributers
}

func (pr *PotRank) GetLoser() []int {

	if len(pr.groups) == 0 {
		return []int{}
	}

	contributers := make([]int, 0)
	for i, g := range pr.groups {
		if i == 0 {
			continue
		}

		contributers = append(contributers, g.Contributers...)
	}

	return contributers
}
