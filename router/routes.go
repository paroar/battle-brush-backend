package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/paroar/battle-brush-backend/game"
)

// CreatePrivateRoom POST Method creates a Room and runs it returning his ID
func CreatePrivateRoom(lobby *game.Lobby, rw http.ResponseWriter, r *http.Request) {
	var _roomJSON RoomJSON

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Couldn't read body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytes, &_roomJSON)
	if err != nil {
		http.Error(rw, "Couldn't Unmarshal", http.StatusInternalServerError)
		return
	}

	roomOptions := &game.RoomOptions{
		NumPlayers: _roomJSON.NumPlayers,
		Time:       _roomJSON.Time,
		Rounds:     _roomJSON.Rounds,
	}

	client, err := lobby.GetClient(_roomJSON.ID)
	if err != nil {
		http.Error(rw, "Client not found", http.StatusNotFound)
		return
	}

	room := lobby.CreatePrivateRoom(roomOptions, client)

	roomJSON := &RoomIDJSON{
		ID: room.ID,
	}

	res, err := json.Marshal(roomJSON)
	if err != nil {
		http.Error(rw, "Couldn't Marshal", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(res)
}

// JoinPrivateRoom PATCH Method joins a Client into a Room
func JoinPrivateRoom(lobby *game.Lobby, rw http.ResponseWriter, req *http.Request) {
	var _joinRoomStruct JoinRoomJSON

	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Couldn't read body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(s, &_joinRoomStruct)
	if err != nil {
		http.Error(rw, "Couldn't Unmarhsal", http.StatusBadRequest)
		return
	}

	client, err := lobby.GetClient(_joinRoomStruct.ClientID)
	if err != nil {
		http.Error(rw, "Client not found", http.StatusNotFound)
		return
	}

	room, err := lobby.GetPrivateRoom(_joinRoomStruct.RoomID)
	if err != nil {
		http.Error(rw, "Room not found", http.StatusConflict)
		return
	}

	err = lobby.JoinPrivateRoom(room, client)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Joined"))
}

// CreateOrJoinRoom POST Method creates a Room and runs it returning his ID
func CreateOrJoinRoom(lobby *game.Lobby, rw http.ResponseWriter, r *http.Request) {
	var _clientIDJSON ClientIDJSON

	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Couldn't read body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(s, &_clientIDJSON)
	if err != nil {
		http.Error(rw, "Couldn't Unmarshal", http.StatusInternalServerError)
		return
	}

	client, err := lobby.GetClient(_clientIDJSON.ID)
	if err != nil {
		http.Error(rw, "Client not found", http.StatusNotFound)
		return
	}

	room := lobby.CreateOrJoinPublicRoom(client)

	roomJSON := &RoomIDJSON{
		ID: room.ID,
	}

	res, err := json.Marshal(roomJSON)
	if err != nil {
		http.Error(rw, "Couldn't Marshal", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(res)
}
