package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/paroar/battle-brush-backend/lobby"
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
