package settlement

import "sort"

type RankGroup struct {
	Contributors []int
	Score        int
}

type PotRank struct {
	contributerCount int
	groups           []*RankGroup
}

func NewPotRank() *PotRank {
	return &PotRank{
		contributerCount: 0,
		groups:           make([]*RankGroup, 0),
	}
}

func (pr *PotRank) AddContributor(score int, contributerIdx int) {

	pr.contributerCount++

	for _, g := range pr.groups {

		if g.Score == score {
			g.Contributors = append(g.Contributors, contributerIdx)
			return
		}
	}

	// Not found so create a new one
	g := &RankGroup{
		Contributors: make([]int, 0),
		Score:        score,
	}

	g.Contributors = append(g.Contributors, contributerIdx)

	pr.groups = append(pr.groups, g)
}

func (pr *PotRank) Calculate() {

	// Sort by score
	sort.Slice(pr.groups, func(i, j int) bool {
		return pr.groups[i].Score > pr.groups[j].Score
	})
}

func (pr *PotRank) ContributorCount() int {
	return pr.contributerCount
}

func (pr *PotRank) GetWinners() []int {

	if len(pr.groups) == 0 {
		return []int{}
	}

	return pr.groups[0].Contributors
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

		contributers = append(contributers, g.Contributors...)
	}

	return contributers
}
