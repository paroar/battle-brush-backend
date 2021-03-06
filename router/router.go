package router

import (
	"flag"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/paroar/battle-brush-backend/websocket"
)

var addr = flag.String("addr", ":8085", "http server address")

// NewRouter creates a mux.router and sets the endpoints
func NewRouter() *http.Server {

	flag.Parse()

	r := mux.NewRouter().StrictSlash(true)

	l := websocket.NewLobby()
	r.Handle("/ws", l)
	r.Handle("/ws/{room}", l)

	r.HandleFunc("/private/{userid}", func(rw http.ResponseWriter, r *http.Request) {
		HandlePrivateRoom(l, rw, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/public/{userid}", func(rw http.ResponseWriter, r *http.Request) {
		HandlePublicRoom(l, rw, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/startgame/{roomid}", func(rw http.ResponseWriter, r *http.Request) {
		HandleStartGame(l, rw, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/chat", func(rw http.ResponseWriter, r *http.Request) {
		HandleChat(l, rw, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/img", HandleImg).Methods(http.MethodPost)
	r.HandleFunc("/vote", HandleVote).Methods(http.MethodPost)

	origins := []string{
		"http://localhost:3000",
	}
	allowedOrigins := handlers.AllowedOrigins(origins)

	handler := handlers.CORS(allowedOrigins)
	s := http.Server{
		Addr:    *addr,
		Handler: handler(r),
	}

	return &s
}
