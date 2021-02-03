package main

import (
	"log"

	"github.com/paroar/battle-brush-backend/router"
)

func main() {
	s := router.NewRouter()
	go log.Fatal(s.ListenAndServe())
}
