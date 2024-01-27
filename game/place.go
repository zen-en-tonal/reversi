package game

import "fmt"

type PlacePiece struct {
	who   Color
	place Place
	flips []flipPiece
}

func placePiece(x int, y int, c Color) PlacePiece {
	var p PlacePiece
	p.place.x = x
	p.place.y = y
	p.who = c
	return p
}

func (c PlacePiece) WhoDoes() Color {
	return c.who
}

func (c PlacePiece) Commit(b *Board) error {
	err := b.placePiece(c.place.x, c.place.y, c.WhoDoes())
	if err != nil {
		return err
	}
	for _, c := range c.flips {
		err = b.MakeEffect(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c PlacePiece) Describe() string {
	return fmt.Sprintf("%s was placed at {x: %d, y: %d}.", c.WhoDoes(), c.place.x, c.place.y)
}

func (c PlacePiece) Score() int {
	return len(c.flips)
}
