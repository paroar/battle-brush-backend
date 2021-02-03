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

	r.HandleFunc("/signin", func(rw http.ResponseWriter, r *http.Request) {
		SignIn(lobby, rw, r)
	}).Methods(http.MethodPost)

	r.HandleFunc("/logout", func(rw http.ResponseWriter, r *http.Request) {
		LogOut(lobby, rw, r)
	}).Methods(http.MethodDelete)

	//ROOMS
	r.HandleFunc("/lobbyrooms", func(rw http.ResponseWriter, r *http.Request) {
		GetLobbyRooms(lobby, rw, r)
	}).Methods(http.MethodGet)

	r.HandleFunc("/createroom", func(rw http.ResponseWriter, r *http.Request) {
		CreateRoom(lobby, rw, r)
	}).Methods(http.MethodPost)

	r.HandleFunc("/joinroom", func(rw http.ResponseWriter, r *http.Request) {
		JoinRoom(lobby, rw, r)
	}).Methods(http.MethodPatch)

	//CLIENTS
	r.HandleFunc("/clients", func(rw http.ResponseWriter, r *http.Request) {
		GetClients(lobby, rw, r)
	}).Methods(http.MethodGet)

	origins := []string{
		"http://localhost:3000",
	}
	allowedOrigins := handlers.AllowedOrigins(origins)
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
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
