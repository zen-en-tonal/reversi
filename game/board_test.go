package game

import (
	"fmt"
	"reflect"
	"testing"
)

func TestHints(t *testing.T) {
	b := NewBoard()

	b.placePiece(3, 2, BLACK)
	b.placePiece(3, 3, BLACK)

	hints, err := b.Hints(WHITE)
	if err != nil {
		t.Error(err)
		return
	}
	if len(hints) != 3 {
		t.Error(hints)
	}
	t.Log(hints)
}

func TestUndo(t *testing.T) {
	b := NewBoard()

	p, _ := b.Hints(WHITE)
	b.MakeEffect(p[0])
	b.Undo(1)

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
		hints, err := board.Hints(c.Color)
		if err != nil {
			t.Error(err)
			return
		}
		var command Command
		for _, h := range hints {
			if h.place == c.Place {
				command = h
				break
			}
		}
		if command == nil {
			t.Error(hints, board)
			return
		}
		err = board.MakeEffect(command)
		if err != nil {
			t.Error(err)
			return
		}
	}

	score := board.Score()
	if score[BLACK] != 8 || score[WHITE] != 1 {
		t.Error("invalid score")
	}

	fmt.Println(board.Logs())
}
