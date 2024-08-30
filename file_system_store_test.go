package poker

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
	{"Name": "Cleo", "Wins": 10},	
	{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemStore(database)

		got := store.league

		want := League{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertNoError(t, err)
		AssertLeague(t, got, want)

		got2 := store.league
		AssertLeague(t, got2, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
	{"Name": "Cleo", "Wins": 10},	
	{"Name": "Chris", "Wins": 33}	
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemStore(database)
		got := store.GetPlayerScore("Chris")
		want := 33

		assertNoError(t, err)
		assertScoreEquals(t, got, want)
	})
	t.Run("store win for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)

		defer cleanDatabase()

		store, err := NewFileSystemStore(database)
		player := "Chris"
		store.RecordWin(player)

		got := store.GetPlayerScore(player)

		assertNoError(t, err)
		assertScoreEquals(t, got, 34)
	})
	t.Run("store win for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemStore(database)
		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		assertNoError(t, err)
		assertScoreEquals(t, got, want)
	})
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemStore(database)

		assertNoError(t, err)
	})
	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)

		defer cleanDatabase()

		store, err := NewFileSystemStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := League{
			{"Chris", 33},
			{"Cleo", 10},
		}

		AssertLeague(t, got, want)

		got = store.GetLeague()
		AssertLeague(t, got, want)
	})
}

// Returns a temp file for persisting our data and the method the will do the garbage collection.
func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}
	return tmpfile, removeFile
}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("incorrect score: got %d, want %d", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("did not expect an error but got one %v", got)
	}

}
