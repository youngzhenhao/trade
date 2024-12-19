package pool

type CandleStick struct {
	Timestamp int `json:"timestamp"`
	Open      int `json:"open"`
	High      int `json:"high"`
	Low       int `json:"low"`
	Close     int `json:"close"`
	Volume    int `json:"volume"`
	Turnover  int `json:"turnover"`
}

type CandleStickCharts struct {
	List []CandleStick `json:"list"`
}
