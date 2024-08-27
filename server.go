package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

type PlayerStore interface {
	GetPlayerScore(string) int
	RecordWin(string)
	GetLeague() []Player
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

type Player struct {
	Name string `json:"name"`
	Wins int    `json:"wins"`
}

// func WriteJSON(w io.Writer, player Player) error {
// 	jsonData, err := json.Marshal(player)
// 	if err != nil {
// 		return err
// 	}
// 	w.Write(jsonData)
// 	return nil
// }

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.GetLeague())
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) GetLeague() []Player {

	return []Player{
		{"Cleo", 32},
		{"Chris", 20},
		{"Tiest", 14},
	}
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	playerName := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, playerName)
	case http.MethodGet:
		p.showScore(w, playerName)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, playerName string) {

	score := p.store.GetPlayerScore(playerName)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)

}

func (p *PlayerServer) processWin(w http.ResponseWriter, playerName string) {
	p.store.RecordWin(playerName)
	w.WriteHeader(http.StatusAccepted)

}

func (s *PlayerServer) RecordWin(name string) {

}
