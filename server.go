package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Player struct {
	Name  string `json: "name"`
	Score int    `json: "score"`
}

var players = []Player{
	{
		Name:  "Pepper",
		Score: 20,
	},
	{
		Name:  "Floyd",
		Score: 15,
	},
}

type PlayerStore interface {
	GetPlayerScore(string) int
	RecordWin(string)
}

type PlayerServer struct {
	store PlayerStore
}

func WriteJSON(w io.Writer, player Player) error {
	jsonData, err := json.Marshal(player)
	if err != nil {
		return err
	}
	w.Write(jsonData)
	return nil
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
