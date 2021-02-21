package lobby

import "errors"

// Rooms struct
type Rooms struct {
	privateRooms map[*Room]bool
	publicRooms  map[*Room]bool
}

// NewRooms returns Rooms
func NewRooms() *Rooms {
	return &Rooms{
		privateRooms: make(map[*Room]bool),
		publicRooms:  make(map[*Room]bool),
	}
}

// GetPublicRoom returns the Room if found or Error
func (rooms *Rooms) GetPublicRoom(id string) (*Room, error) {
	for room := range rooms.publicRooms {
		if room.ID == id {
			return room, nil
		}
	}
	return nil, errors.New("Room not found")
}

// GetPublicRooms returns the Room if found or Error
func (rooms *Rooms) GetPublicRooms() []*Room {
	publicRooms := []*Room{}
	for room := range rooms.publicRooms {
		publicRooms = append(publicRooms, room)
	}
	return publicRooms
}

// GetPrivateRoom returns the Room if found or Error
func (rooms *Rooms) GetPrivateRoom(id string) (*Room, error) {
	for room := range rooms.privateRooms {
		if room.ID == id {
			return room, nil
		}
	}
	return nil, errors.New("Room not found")
}

// GetPrivateRooms returns the Room if found or Error
func (rooms *Rooms) GetPrivateRooms() []*Room {
	privateRooms := []*Room{}
	for room := range rooms.privateRooms {
		privateRooms = append(privateRooms, room)
	}
	return privateRooms
}
