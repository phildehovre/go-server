package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	poker "github.com/phildehovre/go-server"
)

func main() {
	fmt.Println("helloe world")

	file, _ := os.Open("game.db.json")
	store, err := poker.NewFileSystemStore(file)
	server := poker.NewPlayerServer(store)

	if err != nil {
		log.Fatal("no file was found")
	}

	log.Fatal(http.ListenAndServe(":5000", server))
}
