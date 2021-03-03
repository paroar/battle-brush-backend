package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paroar/battle-brush-backend/db"
	"github.com/paroar/battle-brush-backend/game"
	"github.com/paroar/battle-brush-backend/lobby"
	"github.com/paroar/battle-brush-backend/model"
	"github.com/paroar/battle-brush-backend/websocket"
)

//PrivateRoomHandler handler for creating private rooms
func PrivateRoomHandler(l *lobby.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	clientid := vars["userid"]

	client, err := l.GetClient(clientid)
	if err != nil {
		http.Error(rw, "Client not found", http.StatusBadRequest)
		return
	}

	room := lobby.NewPrivateRoom().(lobby.IRoom)
	l.AddRoom(room)
	room.JoinClient(client)

	var rJSON RoomIDJSON
	rJSON.ID = room.GetID()
	res, _ := json.Marshal(&rJSON)

	rw.WriteHeader(http.StatusOK)
	rw.Write(res)

}

//PublicRoomHandler handler for creating public rooms
func PublicRoomHandler(l *lobby.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	clientid := vars["userid"]

	client, err := l.GetClient(clientid)
	if err != nil {
		http.Error(rw, "Client not found", http.StatusBadRequest)
		return
	}

	ro := l.FirstAvailablePublicRoom()
	var room lobby.IRoom
	if ro == nil {
		room = lobby.NewPublicRoom().(lobby.IRoom)
		l.AddRoom(room)
	} else {
		room = ro.(lobby.IRoom)
	}

	room.JoinClient(client)

	var rJSON RoomIDJSON
	rJSON.ID = room.GetID()
	res, _ := json.Marshal(&rJSON)

	rw.WriteHeader(http.StatusOK)
	rw.Write(res)

}

//StartGameHandler handler for starting games
func StartGameHandler(l *lobby.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	roomid := vars["roomid"]

	room := l.GetRoom(roomid)
	var ro lobby.IRoom
	if room != nil {
		ro = room.(lobby.IRoom)
	} else {
		http.Error(rw, "Room not found", http.StatusBadRequest)
		return
	}

	game := ro.GetGame().(lobby.IGame)
	go game.StartGame()

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(""))

}

//HandlePrivateRoom handler for creating private rooms
func HandlePrivateRoom(l *websocket.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	clientid := vars["userid"]

	_, err := db.ReadPlayer(clientid)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	room := model.NewRoom()
	db.CreateRoom(room)
	updatedPlayers := append(room.PlayersID, clientid)
	room = model.UpdateRoom(room.ID, updatedPlayers, room.State)
	db.UpdateRoom(room)
	db.AddPrivateRoom(room.ID)

	var rJSON RoomIDJSON
	rJSON.ID = room.ID
	res, _ := json.Marshal(&rJSON)

	rw.WriteHeader(http.StatusOK)
	rw.Write(res)
}

//HandlePublicRoom handler for creating public rooms
func HandlePublicRoom(l *websocket.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	clientid := vars["userid"]

	_, err := db.ReadPlayer(clientid)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	roomid, err := db.AvailablePublicRoom()
	var room *model.Room
	if err != nil {
		room = model.NewRoom()
		db.CreateRoom(room)
		db.AddPublicRoom(room.ID)
	} else {
		room = db.ReadRoom(roomid)
		updatedPlayers := append(room.PlayersID, clientid)
		updatedRoom := model.UpdateRoom(roomid, updatedPlayers, room.State)
		db.UpdateRoom(updatedRoom)
	}

	var rJSON RoomIDJSON
	rJSON.ID = room.ID
	res, _ := json.Marshal(&rJSON)

	rw.WriteHeader(http.StatusOK)
	rw.Write(res)
}

//HandleStartGame handler for starting games
func HandleStartGame(l *websocket.Lobby, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	roomid := vars["roomid"]

	room := db.ReadRoom(roomid)
	room.State = "Drawing"
	db.CreateRoom(room)

	game := game.NewDrawGame(roomid, room.PlayersID)
	go game.StartGame()

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(""))
}
