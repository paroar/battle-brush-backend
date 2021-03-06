package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paroar/battle-brush-backend/db"
	"github.com/paroar/battle-brush-backend/games"
	"github.com/paroar/battle-brush-backend/message"
	"github.com/paroar/battle-brush-backend/message/content"
	"github.com/paroar/battle-brush-backend/model"
	"github.com/paroar/battle-brush-backend/websocket"
)

//HandlePrivateRoom handler for creating private rooms
func HandlePrivateRoom(l *websocket.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	clientid := vars["userid"]

	player, err := db.ReadPlayer(clientid)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	room := model.NewRoom(clientid, "Private")
	db.CreateRoom(room)

	playersNames := db.ReadPlayersNames(room.PlayersID)
	msg := &message.Envelope{
		Type: content.TypePlayers,
		Content: content.Players{
			UserNames: playersNames,
		},
	}
	l.Broadcast(room.PlayersID, msg)

	msg = &message.Envelope{
		Type: content.TypeJoinLeave,
		Content: content.JoinLeave{
			UserName: player.Name,
			ID:       player.ID,
			Msg:      fmt.Sprintf("%s has joined", player.Name),
		},
	}
	l.Broadcast(room.PlayersID, msg)

	var rJSON roomIDJSON
	rJSON.ID = room.ID
	res, _ := json.Marshal(&rJSON)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res)
}

//HandlePublicRoom handler for creating public rooms
func HandlePublicRoom(l *websocket.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	clientid := vars["userid"]

	player, err := db.ReadPlayer(clientid)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	room, err := db.AvailablePublicRoom()
	if err != nil {
		room = model.NewRoom(clientid, "Public")
		db.CreateRoom(room)
	} else {
		updatedPlayers := append(room.PlayersID, clientid)
		room.UpdateRoom(updatedPlayers, room.State)
		db.UpdateRoom(room)
	}

	player.RoomID = room.ID
	db.UpdatePlayer(player)

	playersNames := db.ReadPlayersNames(room.PlayersID)
	msg := &message.Envelope{
		Type: content.TypePlayers,
		Content: content.Players{
			UserNames: playersNames,
		},
	}
	l.Broadcast(room.PlayersID, msg)

	msg = &message.Envelope{
		Type: content.TypeJoinLeave,
		Content: content.JoinLeave{
			UserName: player.Name,
			ID:       player.ID,
			Msg:      fmt.Sprintf("%s has joined", player.Name),
		},
	}
	l.Broadcast(room.PlayersID, msg)

	var rJSON roomIDJSON
	rJSON.ID = room.ID
	res, _ := json.Marshal(&rJSON)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(res)
}

//HandleStartGame handler for starting games
func HandleStartGame(l *websocket.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	roomid := vars["roomid"]

	room, err := db.ReadRoom(roomid)
	if err != nil {
		http.Error(rw, "Couldn't start the game", http.StatusBadRequest)
		return
	}
	room.State = "Drawing"
	db.UpdateRoom(room)

	game := games.NewDrawGame(roomid, room.PlayersID)
	go game.StartGame(l)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte("Game started"))
}

//HandleChat handler manages chat
func HandleChat(l *websocket.Lobby, rw http.ResponseWriter, r *http.Request) {
	var chat chat
	json.NewDecoder(r.Body).Decode(&chat)

	room, err := db.ReadRoom(chat.Roomid)
	if err != nil {
		http.Error(rw, "Room not found", http.StatusBadRequest)
		return
	}

	msg := &message.Envelope{
		Type:    content.TypeChat,
		Content: chat,
	}

	l.Broadcast(room.PlayersID, msg)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(chat)
}

//HandleImg handler manages img
func HandleImg(rw http.ResponseWriter, r *http.Request) {
	var img img
	json.NewDecoder(r.Body).Decode(&img)

	drawing := model.NewDrawing(img.Playerid, img.Playerid, img.Img)
	db.CreateDrawing(drawing)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(img)
}

//HandleVote handler manages img
func HandleVote(rw http.ResponseWriter, r *http.Request) {
	var v vote
	json.NewDecoder(r.Body).Decode(&v)

	vote := model.NewVote(v.PlayerID, v.Vote)
	db.CreateVote(vote)

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(vote)
}
