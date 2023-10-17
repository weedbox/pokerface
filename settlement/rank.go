package settlement

import "sort"

type RankGroup struct {
	Contributors []int
	Score        int
}

type Rank struct {
	contributerCount int
	groups           []*RankGroup
}

func NewRank() *Rank {
	return &Rank{
		contributerCount: 0,
		groups:           make([]*RankGroup, 0),
	}
}

func (r *Rank) AddContributor(score int, contributerIdx int) {

	r.contributerCount++

	for _, g := range r.groups {

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

	r.groups = append(r.groups, g)
}

func (r *Rank) Calculate() {

	// Sort by score
	sort.Slice(r.groups, func(i, j int) bool {
		return r.groups[i].Score > r.groups[j].Score
	})
}

func (r *Rank) ContributorCount() int {
	return r.contributerCount
}

func (r *Rank) GetWinners() []int {

	if len(r.groups) == 0 {
		return []int{}
	}

	return r.groups[0].Contributors
}

func (r *Rank) GetLoser() []int {

	if len(r.groups) == 0 {
		return []int{}
	}

	contributers := make([]int, 0)
	for i, g := range r.groups {
		if i == 0 {
			continue
		}

		contributers = append(contributers, g.Contributors...)
	}

	return contributers
}
