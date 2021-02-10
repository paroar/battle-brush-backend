package router

import (
	"flag"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/paroar/battle-brush-backend/game"
)

var addr = flag.String("addr", ":8085", "http server address")

// NewRouter creates a mux.router and sets the endpoints
func NewRouter() *http.Server {

	flag.Parse()

	r := mux.NewRouter()

	lobby := game.NewLobby()
	go lobby.Run()

	//LOBBY
	r.Handle("/ws", lobby)

	//ROOMS
	r.HandleFunc("/createroom", func(rw http.ResponseWriter, r *http.Request) {
		CreatePrivateRoom(lobby, rw, r)
	}).Methods(http.MethodPost)

	r.HandleFunc("/joinroom", func(rw http.ResponseWriter, r *http.Request) {
		JoinPrivateRoom(lobby, rw, r)
	}).Methods(http.MethodPatch)

	r.HandleFunc("/createjoin", func(rw http.ResponseWriter, r *http.Request) {
		CreateOrJoinRoom(lobby, rw, r)
	}).Methods(http.MethodPost)

	origins := []string{
		"http://localhost:3000",
	}
	allowedOrigins := handlers.AllowedOrigins(origins)

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodOptions,
	}
	allowedMethods := handlers.AllowedMethods(methods)

	handler := handlers.CORS(allowedOrigins, allowedMethods)
	s := http.Server{
		Addr:    *addr,
		Handler: handler(r),
	}

	return &s
}
