package router

type roomIDJSON struct {
	ID string `json:"roomid"`
}

type chat struct {
	Roomid   string `json:"roomid"`
	Playerid string `json:"playerid"`
	Username string `json:"username"`
	Msg      string `json:"msg"`
}

type img struct {
	Playerid string `json:"playerid"`
	Roomid   string `json:"roomid"`
	Img      string `json:"img"`
}

type vote struct {
	PlayerID string  `json:"playerid"`
	Vote     float64 `json:"vote"`
}
