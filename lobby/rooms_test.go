package lobby

import "testing"

func TestGetPublicRoom(t *testing.T) {
	l := NewLobby()
	go l.Run()

	room := NewPublicRoom()

	go room.run()

	l.joinPublicRoom(room)

	_, err := l.GetPublicRoom(room.ID)

	if err != nil {
		t.Fatal("GetPublicRoom should return the room ID")
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

func TestGetPublicRooms(t *testing.T) {
	l := NewLobby()
	go l.Run()

	room := NewPublicRoom()

	go room.run()

	l.joinPublicRoom(room)

	publicRooms := l.rooms.GetPublicRooms()

	if len(publicRooms) != 1 {
		t.Fatal("GetPublicRooms should return an slice of rooms")
	}
}

func TestGetPrivateRooms(t *testing.T) {
	l := NewLobby()
	go l.Run()

	room := NewPrivateRoom(l)

	go room.run()

	l.joinPrivateRoom(room)

	privateRooms := l.rooms.GetPrivateRooms()

	if len(privateRooms) != 1 {
		t.Fatal("GetPrivateRooms should return an slice of rooms")
	}
}
