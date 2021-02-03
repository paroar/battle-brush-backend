package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/paroar/battle-brush-backend/game"
)

// SignIn POST Method that creates a new Client and returns the userID
func SignIn(lobby *game.Lobby, rw http.ResponseWriter, r *http.Request) {
	client := game.NewClient(lobby)

	lobby.JoinClientChan <- client

	var _clientID = &ClientIDJSON{
		ID: client.ID,
	}

	response, err := json.Marshal(_clientID)
	if err != nil {
		http.Error(rw, "Couldn't Marshal", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

// LogOut DELETE Method creates a Room and runs it returning his ID
func LogOut(lobby *game.Lobby, rw http.ResponseWriter, req *http.Request) {
	var _clientIDJSON ClientIDJSON

	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Couldn't read body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(s, &_clientIDJSON)
	if err != nil {
		http.Error(rw, "Couldn't Unmarshal", http.StatusInternalServerError)
		return
	}

	client, err := lobby.GetLobbyClient(_clientIDJSON.ID)
	if err != nil {
		http.Error(rw, "Client not found", http.StatusNotFound)
		return
	}

	lobby.LeaveClientChan <- client

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

}

// GetLobbyRooms GET Method returns all Rooms in Lobby
func GetLobbyRooms(lobby *game.Lobby, rw http.ResponseWriter, r *http.Request) {
	var lrs []LobbyRoomJSON

	rooms := lobby.GetLobbyRooms()
	for _, room := range rooms {
		clients := []string{}
		for client := range room.Clients {
			clients = append(clients, client.ID)
		}
		lr := LobbyRoomJSON{
			ID:      room.ID,
			Clients: clients,
		}
		lrs = append(lrs, lr)
	}

	response, err := json.Marshal(lrs)
	if err != nil {
		http.Error(rw, "Couldn't Marshal", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

// CreateRoom POST Method creates a Room and runs it returning his ID
func CreateRoom(lobby *game.Lobby, rw http.ResponseWriter, req *http.Request) {
	var _clientIDJSON ClientIDJSON

	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Couldn't read body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(s, &_clientIDJSON)
	if err != nil {
		http.Error(rw, "Couldn't Unmarshal", http.StatusInternalServerError)
		return
	}

	client, err := lobby.GetLobbyClient(_clientIDJSON.ID)
	if err != nil {
		http.Error(rw, "Client not found", http.StatusNotFound)
		return
	}

	room := game.NewRoom(lobby)

	r := RoomIDJSON{
		ID: room.ID,
	}

	response, err := json.Marshal(r)
	if err != nil {
		http.Error(rw, "Couldn't Marshal", http.StatusInternalServerError)
		return
	}

	go room.Run()
	lobby.JoinRoomChan <- room

	lobby.LeaveClientChan <- client
	room.JoinClientChan <- client

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

// JoinRoom PATCH Method joins a Client into a Room
func JoinRoom(lobby *game.Lobby, rw http.ResponseWriter, req *http.Request) {
	var _joinRoomStruct JoinRoomJSON

	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Couldn't read body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(s, &_joinRoomStruct)
	if err != nil {
		http.Error(rw, "Couldn't Unmarhsal", http.StatusInternalServerError)
		return
	}

	client, err := lobby.GetLobbyClient(_joinRoomStruct.ClientID)
	if err != nil {
		http.Error(rw, "Client not found", http.StatusNotFound)
		return
	}

	room, err := lobby.GetLobbyRoom(_joinRoomStruct.RoomID)
	if err != nil {
		http.Error(rw, "Room not found", http.StatusNotFound)
		return
	}

	lobby.LeaveClientChan <- client
	room.JoinClientChan <- client

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(""))
}

// GetClients GET method to get all Clients in the Lobby
func GetClients(l *game.Lobby, rw http.ResponseWriter, r *http.Request) {
	clients := []string{}
	for _, client := range l.GetLobbyClients() {
		clients = append(clients, client.ID)
	}

	var _clientsJSON = ClientsIDJSON{
		ClientsID: clients,
	}

	response, err := json.Marshal(_clientsJSON)
	if err != nil {
		http.Error(rw, "Couldn't Marshal", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}
