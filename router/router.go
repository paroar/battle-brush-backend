package router

import (
	"flag"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/paroar/battle-brush-backend/lobby"
)

var addr = flag.String("addr", ":8085", "http server address")

// NewRouter creates a mux.router and sets the endpoints
func NewRouter() *http.Server {

	flag.Parse()

	r := mux.NewRouter().StrictSlash(true)

	l := lobby.NewLobby()
	// go l.Run()

	//LOBBY
	r.Handle("/ws", l)
	r.Handle("/ws/{room}", l)

	origins := []string{
		"http://localhost:3000",
	}
	allowedOrigins := handlers.AllowedOrigins(origins)

	methods := []string{
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
