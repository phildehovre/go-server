package poker_test

import (
	"strings"
	"testing"

	poker "github.com/phildehovre/go-server"
)

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input,", func(t *testing.T) {

		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()
		winner := "Chris"

		poker.AssertPlayerWin(t, playerStore, winner)

	})
	t.Run("record cleo win from user input, ", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})
}
