package lobby

import (
	"testing"

	"github.com/google/uuid"
	"github.com/paroar/battle-brush-backend/generators"
)

func TestJoinClient(t *testing.T) {

	l := NewLobby()
	go l.Run()

	client := &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: l,
		send:  make(chan *Message),
	}

	l.joinClient(client)

	if _, err := l.GetClient(client.id); err != nil {
		t.Fatal("Client should be on lobby")
	}

}

func TestLeaveClient(t *testing.T) {

	l := NewLobby()
	go l.Run()

	client := &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: l,
		send:  make(chan *Message),
	}

	l.joinClient(client)
	l.leaveClient(client)

	if _, err := l.GetClient(client.id); err == nil {
		t.Fatal("Client shouldn't be on lobby")
	}

}

func TestGetClient_Client(t *testing.T) {

	l := NewLobby()
	go l.Run()

	client := &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: l,
		send:  make(chan *Message),
	}

	l.joinClient(client)

	if c, _ := l.GetClient(client.id); c != client {
		t.Fatalf("Client should be the same: %+v %+v", c, client)
	}

}

func TestGetClient_NoClient(t *testing.T) {

	l := NewLobby()
	go l.Run()

	client := &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: l,
		send:  make(chan *Message),
	}

	if _, err := l.GetClient(client.id); err == nil {
		t.Fatalf("Shouldn't get the client")
	}

}

func TestFirstAvailablePublicRoom_Room(t *testing.T) {

	l := NewLobby()
	go l.Run()

	room := NewDefaultRoom(l)

	l.joinPublicRoom(room)

	availableRoom := l.firstAvailablePublicRoom()

	if availableRoom == nil {
		t.Fatal("Should be a room")
	}

}

func TestFirstAvailablePublicRoom_NoRoom(t *testing.T) {

	l := NewLobby()
	go l.Run()

	room := l.firstAvailablePublicRoom()

	if room != nil {
		t.Fatal("Shouldn't be any rooms")
	}

}

func TestCreateOrJoinPublicRoom(t *testing.T) {
	l := NewLobby()
	go l.Run()

	client := &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: l,
		send:  make(chan *Message),
	}

	r := l.CreateOrJoinPublicRoom(client)
	if r == nil {
		t.Fatal("CreateOrJoinPublicRoom should return a room")
	}
}

func TestLeavePublicRoom(t *testing.T) {
	l := NewLobby()
	go l.Run()

	room := NewDefaultRoom(l)

	go room.run()

	l.joinPublicRoom(room)
	l.leavePublicRoom(room)

	_, err := l.GetPublicRoom(room.ID)
	if err == nil {
		t.Fatal("rooms.GetPublicRoom should return an error because there is no room")
	}

}

func TestJoinPrivateRoom(t *testing.T) {
	l := NewLobby()
	go l.Run()

	room := NewPrivateRoom(l)

	go room.run()

	l.joinPrivateRoom(room)

	client := &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: l,
		send:  make(chan *Message),
	}

	err := l.JoinClientToPrivateRoom(room, client)

	if err != nil {
		t.Fatal("JoinPrivateRoom should add the client into the private room")
	}

}

func TestLeavePrivateRoom(t *testing.T) {
	l := NewLobby()
	go l.Run()

	room := NewPrivateRoom(l)

	go room.run()

	l.joinPrivateRoom(room)
	l.leavePrivateRoom(room)

	_, err := l.GetPrivateRoom(room.ID)
	if err == nil {
		t.Fatal("GetPrivateRoom should return an error because there is no room")
	}

}

func TestJoinPrivateRoom_Full(t *testing.T) {
	l := NewLobby()
	go l.Run()

	clients := make(map[*Client]bool)
	roomOptions := RoomOptions{
		NumPlayers: 0,
	}
	room := &Room{
		lobby:           l,
		clients:         clients,
		joinClientChan:  make(chan *Client),
		leaveClientChan: make(chan *Client),
		ID:              uuid.NewString(),
		broadcast:       make(chan *Message),
		options:         roomOptions,
		game:            NewDrawGame(clients),
	}

	go room.run()

	l.joinPrivateRoom(room)

	client := &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: l,
		send:  make(chan *Message),
	}

	err := l.JoinClientToPrivateRoom(room, client)

	if err == nil {
		t.Fatal("JoinPrivateRoom should return an error")
	}

}

func TestGetPrivateRoom(t *testing.T) {
	l := NewLobby()
	go l.Run()

	room := NewPrivateRoom(l)

	go room.run()

	l.joinPrivateRoom(room)

	_, err := l.GetPrivateRoom(room.ID)

	if err != nil {
		t.Fatal("GetPrivateRoom should return the room ID")
	}
}

func TestCreatePrivateRoom(t *testing.T) {
	l := NewLobby()
	go l.Run()

	client := &Client{
		name:  generators.Name(),
		id:    uuid.NewString(),
		lobby: l,
		send:  make(chan *Message),
	}

	r := l.CreatePrivateRoom(client)
	if r == nil {
		t.Fatal("CreatePrivateRoom should return the created room ID")
	}
}

func TestAvailableRooms(t *testing.T) {
	l := NewLobby()
	go l.Run()

	clients := make(map[*Client]bool)

	room1 := &Room{
		lobby:           l,
		clients:         clients,
		joinClientChan:  make(chan *Client),
		leaveClientChan: make(chan *Client),
		ID:              uuid.NewString(),
		broadcast:       make(chan *Message),
		options: RoomOptions{
			NumPlayers: 0,
		},
		game: NewDrawGame(clients),
	}

	l.joinPrivateRoom(room1)

	room2 := &Room{
		lobby:           l,
		clients:         clients,
		joinClientChan:  make(chan *Client),
		leaveClientChan: make(chan *Client),
		ID:              uuid.NewString(),
		broadcast:       make(chan *Message),
		options: RoomOptions{
			NumPlayers: 1,
		},
		game: NewDrawGame(clients),
	}

	l.joinPrivateRoom(room2)

	availablePrivateRooms := l.availableRooms(l.rooms.GetPrivateRooms())

	roomFound := false
	for _, r := range availablePrivateRooms {
		if r == room2 {
			roomFound = true
		}
	}

	if !roomFound {
		t.Fatal("availableRooms should return the available rooms")
	}

}
