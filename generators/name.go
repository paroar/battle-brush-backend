package generators

import (
	"fmt"
	"math/rand"
	"time"
)

// Name returns a random name
func Name() string {
	rand.Seed(time.Now().Unix())
	adj := adj[rand.Intn(len(adj))]
	name := names[rand.Intn(len(names))]
	fullName := fmt.Sprintf("%s_%s", adj, name)
	return fullName
}

var adj = []string{
	"bad",
	"best",
	"better",
	"big",
	"black",
	"certain",
	"clear",
	"different",
	"early",
	"easy",
	"economic",
	"federal",
	"free",
	"full",
	"good",
	"great",
	"hard",
	"high",
	"human",
	"important",
	"international",
	"large",
	"late",
	"little",
	"local",
	"long",
	"low",
	"major",
	"military",
	"national",
	"new",
	"old",
	"only",
	"other",
	"political",
	"possible",
	"public",
	"real",
	"recent",
	"right",
	"small",
	"social",
	"special",
	"strong",
	"sure",
	"true",
	"white",
	"whole",
	"young",
}

var names = []string{
	"Donatello",
	"Botticelli",
	"Leonardo",
	"Michelangelo",
	"Raphael",
	"Titian",
	"Durer",
	"Greco",
	"Caravaggio",
	"Bernini",
	"Rembrandt",
	"Goya",
	"Pissarro",
	"Manet",
	"Degas",
	"Matisse",
	"Klee",
	"Picasso",
	"Hopper",
	"Kandinsky",
	"van",
	"Miro",
	"Rockwell",
	"Magritte",
	"Escher",
	"Dali",
	"Pollock",
	"Warhol",
	"Giotto",
	"Velazquez",
	"Masaccio",
	"Courbet",
}
