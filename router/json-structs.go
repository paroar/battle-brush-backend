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
