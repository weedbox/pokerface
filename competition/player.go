package competition

type PlayerInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Bankroll     int64  `json:"bankroll"`
	Participated bool   `json:"participated"`
}
