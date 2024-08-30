package poker

import (
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
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {

	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league

}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %s", len(store.winCalls), winner)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}

func NewGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func AssertStatus(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("incorrect status code: got %d want %d", got, want)
	}
}
func AssertLeague(t *testing.T, got, wantedLeague []Player) {
	if !reflect.DeepEqual(got, wantedLeague) {
		t.Errorf("got %v want %v", got, wantedLeague)
	}
}

func AssertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	got := response.Result().Header.Get("content-type")
	if got != want {
		t.Errorf("response did not have content-type of %s, got %v", want, got)
	}
}

func NewLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func GetLeagueFromResponse(t testing.TB, body io.Reader) []Player {
	t.Helper()
	league, _ := NewLeague(body)

	return league
}
