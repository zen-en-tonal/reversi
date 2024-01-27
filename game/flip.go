package game

import "errors"

type FlipPiece struct {
	place Place
}

func (c FlipPiece) WhoDoes() Color {
	return None
}

func (c FlipPiece) Commit(b *Board) error {
	color, err := b.GetPiece(c.place.x, c.place.y)
	if err != nil {
		return err
	}
	if *color == None {
		return errors.New("invalid operation")
	}
	return b.placePiece(c.place.x, c.place.y, color.Opposite())
}

func (c FlipPiece) Describe() string {
	return "flips."
}
