package game

import "errors"

type flipPiece struct {
	place Place
}

func (c flipPiece) WhoDoes() Color {
	return None
}

func (c flipPiece) Commit(b *Board) error {
	color, err := b.GetPiece(c.place)
	if err != nil {
		return err
	}
	if *color == None {
		return errors.New("invalid operation")
	}
	return b.placePiece(c.place, color.Opposite())
}

func (c flipPiece) Describe() string {
	return "flips."
}
