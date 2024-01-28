package game

import (
	"reflect"
	"testing"
)

func TestHints(t *testing.T) {
	b := NewBoard()

	b.placePiece(NewPlace(3, 2), BLACK)
	b.placePiece(NewPlace(3, 3), BLACK)

	hints := b.Hints(WHITE)
	if len(hints) != 3 {
		t.Errorf("length of hints must be 3 but it was %d", len(hints))
	}
}

func TestUndo(t *testing.T) {
	b := NewBoard()

	p := b.Hints(WHITE)
	b.MustMakeEffect(p[0])
	b.MustUndo(1)

	if !reflect.DeepEqual(b.pieces, NewBoard().pieces) {
		t.Error("invalid states.")
	}
}

func TestGame(t *testing.T) {
	board := NewBoard()

	comamnds := []struct {
		Place
		Color
	}{
		{Place{x: 3, y: 2}, BLACK},
		{Place{x: 2, y: 2}, WHITE},
		{Place{x: 4, y: 5}, BLACK},
		{Place{x: 3, y: 5}, WHITE},
		{Place{x: 3, y: 6}, BLACK},
	}

	for _, c := range comamnds {

		hints := board.Hints(c.Color)
		var command Command
		for _, h := range hints {
			if h.place == c.Place {
				command = h
				break
			}
		}
		if command == nil {
			t.Errorf("expected command is not found in hints.")
			return
		}

		board.MustMakeEffect(command)
	}

	score := board.Score()
	if score[BLACK] != 8 || score[WHITE] != 1 {
		t.Errorf("expected score is {black: 8, white: 1} but actual is {black: %d, white: %d}", score[BLACK], score[WHITE])
	}
}
