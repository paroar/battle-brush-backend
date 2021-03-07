package router

type roomJSON struct {
	ID string `json:"roomid"`
}

type chatJSON struct {
	Roomid   string `json:"roomid"`
	Playerid string `json:"playerid"`
	Username string `json:"username"`
	Msg      string `json:"msg"`
}

type imgJSON struct {
	Playerid string `json:"playerid"`
	Roomid   string `json:"roomid"`
	Img      string `json:"img"`
}

type voteJSON struct {
	PlayerID string  `json:"playerid"`
	Vote     float64 `json:"vote"`
}
