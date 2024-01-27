package game

import "fmt"

type Skip struct {
	who Color
}

func (c Skip) WhoDoes() Color {
	return c.who
}

func (c Skip) Commit(_ *Board) error {
	// no-op
	return nil
}

func (c Skip) Describe() string {
	return fmt.Sprintf("%s was skipped.", c.WhoDoes())
}
