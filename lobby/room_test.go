package lobby

// func TestRoomLeaveClient(t *testing.T) {
// 	l := NewLobby()
// 	go l.Run()

// 	client := &Client{
// 		name:  generators.Name(),
// 		id:    uuid.NewString(),
// 		lobby: l,
// 		send:  make(chan *Message),
// 	}

// 	l.joinClient(client)

// 	room := l.CreatePrivateRoom(client)

// 	room.leaveClient(client)

// 	_, err := room.getClient(client.id)
// 	if err == nil {
// 		t.Fatal("room.getClient should return an error")
// 	}

// }

// func TestRoomGetClient(t *testing.T) {
// 	l := NewLobby()
// 	go l.Run()

// 	client := &Client{
// 		name:  generators.Name(),
// 		id:    uuid.NewString(),
// 		lobby: l,
// 		send:  make(chan *Message),
// 	}

// 	l.joinClient(client)

// 	room := l.CreatePrivateRoom(client)

// 	c, err := room.getClient(client.id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if c.id != client.id {
// 		t.Fatal("room.getClient should return the same client")
// 	}
// }
