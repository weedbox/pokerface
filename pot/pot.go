package pot

type Pot struct {
	Wager        int64 `json:"wager"`
	Total        int64 `json:"total"`
	Contributors []int `json:"contributors"`
}
