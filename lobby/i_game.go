package lobby

import "github.com/paroar/battle-brush-backend/drawing"

//IGame game interface
type IGame interface {
	StartGame()
	Vote(*Vote)
	Drawing(*drawing.Drawing)
}
