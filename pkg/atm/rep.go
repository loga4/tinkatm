package atm

type ATM struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Limit    int    `json:"limit"`
	Index    int    `json:"index"`
	Location Point  `json:"location"`
}

type Point struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}
