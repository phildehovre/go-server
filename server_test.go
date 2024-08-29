package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return []Player{}
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		[]string{},
		[]Player{
			{"chris", 30},
		},
	}

	server := NewPlayerServer(&store)
	t.Run("returns Pepper's score", func(t *testing.T) {

		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertResponseBody(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("return Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "10")
		assertStatus(t, response.Code, http.StatusOK)

	})

	t.Run("retuns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func newGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("incorrect status code: got %d want %d", got, want)
	}
}

func TestStoreWins(t *testing.T) {

	t.Run("it returns accepted on POST", func(t *testing.T) {
		store := StubPlayerStore{
			map[string]int{},
			[]string{},
			[]Player{
				{"chris", 30},
			},
		}
		server := NewPlayerServer(&store)
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		if len(store.winCalls) != 1 {
			t.Errorf("want 1, got %d", len(store.winCalls))
		}
		assertStatus(t, response.Code, http.StatusAccepted)
	})
	t.Run("it records wins on POST", func(t *testing.T) {
		store := StubPlayerStore{
			map[string]int{},
			[]string{},
			[]Player{
				{"chris", 30},
			},
		}
		server := NewPlayerServer(&store)
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Fatalf("got %d calls to recordwin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("want %s got %s", player, store.winCalls[0])
		}
	})
}

func newPostWinRequest(playerName string) *http.Request {
	return httptest.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", playerName), nil)
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("unable to parse response from server %q into slice of player '%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)

	})

	t.Run("it returns the league table as json", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)

	})
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) []Player {
	t.Helper()
	league, _ := NewLeague(body)

	return league
}

func assertLeague(t *testing.T, got, wantedLeague []Player) {
	if !reflect.DeepEqual(got, wantedLeague) {
		t.Errorf("got %v want %v", got, wantedLeague)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	got := response.Result().Header.Get("content-type")
	if got != want {
		t.Errorf("response did not have content-type of %s, got %v", want, got)
	}
}
