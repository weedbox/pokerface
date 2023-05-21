package pot

type Pot struct {
	Wager        int64 `json:"wager"`
	Total        int64 `json:"total"`
	Contributers []int `json:"contributers"`
}
