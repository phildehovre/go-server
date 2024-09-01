package poker_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	poker "github.com/phildehovre/go-server"
)

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}
type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(At time.Duration, Amount int) {
	s.alerts = append(s.alerts, ScheduledAlert{At, Amount})
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips At %v", s.Amount, s.At)
}

type GameSpy struct {
	StartedWith  int
	FinishedWith string
	StartCalled  bool
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
}
func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

var DummySpyAlerter = &SpyBlindAlerter{}
var DummyPlayerStore = &poker.StubPlayerStore{}
var DummyStdIn = &bytes.Buffer{}
var DummyStdOut = &bytes.Buffer{}

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input,", func(t *testing.T) {

		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := &GameSpy{}

		cli := poker.NewCLI(in, DummyStdOut, game)
		cli.PlayPoker()
		winner := "Chris"

		poker.AssertPlayerWin(t, playerStore, winner)

	})
	t.Run("record cleo win from user input, ", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := &GameSpy{}

		cli := poker.NewCLI(in, DummyStdOut, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})
	// t.Run("it schedules printing of blind values", func(t *testing.T) {
	// 	in := strings.NewReader("Chris wins\n")
	// 	blindAlerter := &SpyBlindAlerter{}

	// 	stdout := &bytes.Buffer{}
	// 	game := &GameSpy{}
	// 	cli := poker.NewCLI(in, stdout, game)
	// 	cli.PlayPoker()

	// 	cases := []ScheduledAlert{
	// 		{0 * time.Second, 100},
	// 		{10 * time.Minute, 200},
	// 		{20 * time.Minute, 300},
	// 		{30 * time.Minute, 400},
	// 		{40 * time.Minute, 500},
	// 		{50 * time.Minute, 600},
	// 		{60 * time.Minute, 800},
	// 		{70 * time.Minute, 1000},
	// 		{80 * time.Minute, 2000},
	// 		{90 * time.Minute, 4000},
	// 		{100 * time.Minute, 8000},
	// 	}

	// 	for i, want := range cases {
	// 		t.Run(fmt.Sprint(want), func(t *testing.T) {

	// 			if len(blindAlerter.alerts) <= i {
	// 				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
	// 			}

	// 			got := blindAlerter.alerts[i]
	// 			assertScheduledAlert(t, got, want)
	// 		})
	// 	}

	// })
	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		blindAlerter := &SpyBlindAlerter{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}
		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		got := stdout.String()
		want := "Please enter the number of players: "

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
		cases := []ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}
		checkSchedulingCases(cases, t, blindAlerter)

	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &GameSpy{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Errorf("game should not have started")
		}
	})
}
func assertScheduledAlert(t testing.TB, got, want ScheduledAlert) {
	t.Helper()

	AmountGot := got.Amount
	if AmountGot != want.Amount {
		t.Errorf("got Amount %d, want %d", AmountGot, want.Amount)
	}

	gotScheduledTime := got.At
	if gotScheduledTime != want.At {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, want.At)
	}
}

func TestTexasHoldem_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(DummyPlayerStore, blindAlerter)

		game.Start(5)

		cases := []ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 10 * time.Minute, Amount: 200},
			{At: 20 * time.Minute, Amount: 300},
			{At: 30 * time.Minute, Amount: 400},
			{At: 40 * time.Minute, Amount: 500},
			{At: 50 * time.Minute, Amount: 600},
			{At: 60 * time.Minute, Amount: 800},
			{At: 70 * time.Minute, Amount: 1000},
			{At: 80 * time.Minute, Amount: 2000},
			{At: 90 * time.Minute, Amount: 4000},
			{At: 100 * time.Minute, Amount: 8000},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(DummyPlayerStore, blindAlerter)

		game.Start(7)

		cases := []ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Minute, Amount: 200},
			{At: 24 * time.Minute, Amount: 300},
			{At: 36 * time.Minute, Amount: 400},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

}

func TestTexasHoldem_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewTexasHoldem(store, DummySpyAlerter)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}

func checkSchedulingCases(cases []ScheduledAlert, t *testing.T, blindAlerter *SpyBlindAlerter) {
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter)
			}
			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
}
