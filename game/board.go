package game

import (
	"errors"
)

const BOARDSIZE int = 8

type Board struct {
	history []Command
	pieces  map[Place]Color
}

func NewBoard() Board {
	var b Board
	b.pieces = make(map[Place]Color)
	for y := 0; y < BOARDSIZE; y++ {
		for x := 0; x < BOARDSIZE; x++ {
			b.placePiece(x, y, None)
		}
	}
	// 初期配置
	b.placePiece(3, 3, WHITE)
	b.placePiece(4, 3, BLACK)
	b.placePiece(3, 4, BLACK)
	b.placePiece(4, 4, WHITE)

	return b
}

func (b *Board) Clone() Board {
	h := make([]Command, len(b.history))
	copy(h, b.history)
	p := make(map[Place]Color)
	for k, v := range b.pieces {
		p[k] = v
	}
	newB := NewBoard()
	newB.history = h
	newB.pieces = p
	return newB
}

func (b Board) Logs() []string {
	var logs []string
	for _, c := range b.history {
		if c.WhoDoes() == None {
			continue
		}
		logs = append(logs, c.Describe())
	}
	return logs
}

func (b Board) Pieces() map[Place]Color {
	return b.Clone().pieces
}

func (b Board) WhoesTurn() Color {
	return b.history[len(b.history)-1].WhoDoes().Opposite()
}

func (b *Board) Undo(num int) error {
	numToCommit := len(b.history) - num
	if numToCommit < 0 {
		return errors.New("out of range")
	}
	newBoard := NewBoard()
	for _, c := range b.history[:numToCommit] {
		err := newBoard.MakeEffect(c)
		if err != nil {
			return err
		}
	}
	*b = newBoard
	return nil
}

func (b Board) Score() map[Color]int {
	m := make(map[Color]int)
	m[BLACK] = 0
	m[WHITE] = 0
	for _, c := range b.pieces {
		if c == BLACK {
			m[BLACK] += 1
		}
		if c == WHITE {
			m[WHITE] += 1
		}
	}
	return m
}

func (b *Board) MakeEffect(c Command) error {
	nb := b.Clone()
	err := c.Commit(&nb)
	if err != nil {
		return err
	}
	*b = nb
	if c.WhoDoes() != None {
		b.history = append(b.history, c)
	}
	return nil
}

func (b Board) Hints(c Color) ([]PlacePiece, error) {
	var ps []PlacePiece
	for place, color := range b.pieces {

		if color != c {
			continue
		}

		directions := []direction{UP, DOWN, LEFT, RIGHT, UP_LEFT, UP_RIGHT, DOWN_LEFT, DOWN_RIGHT}
		for _, d := range directions {
			path := path{
				current: place.next(d),
			}
			if res := b.scan(c.Opposite(), d, path); res != nil {
				p := PlacePiece{
					who:   c,
					place: res.current,
					flips: res.intoFlips(),
				}
				ps = append(ps, p)
			}
		}

	}

	return ps, nil
}

func (b Board) scan(expect Color, d direction, p path) *path {
	currentColor, err := b.GetPiece(p.current.x, p.current.y)
	if err != nil {
		return nil
	}
	if *currentColor == expect {
		return b.scan(expect, d, p.walkNext(d))
	}
	if *currentColor == None && len(p.log) > 0 {
		return &p
	}
	return nil
}

type path struct {
	current Place
	log     []Place
}

func (p path) walkNext(d direction) path {
	p.log = append(p.log, p.current)
	p.current = p.current.next(d)
	return p
}

func (p path) intoFlips() []FlipPiece {
	var f []FlipPiece
	for _, l := range p.log {
		f = append(f, FlipPiece{place: l})
	}
	return f
}

func (b Board) inBound(p Place) bool {
	return p.x >= 0 || p.x < BOARDSIZE || p.y >= 0 || p.y < BOARDSIZE
}

func (b *Board) placePiece(x int, y int, c Color) error {
	place := Place{x, y}
	if !b.inBound(place) {
		return errors.New("out of range")
	}
	b.pieces[place] = c
	return nil
}

func (b Board) GetPiece(x int, y int) (*Color, error) {
	place := Place{x, y}
	if !b.inBound(place) {
		return nil, errors.New("out of range")
	}
	c := b.pieces[place]
	return &c, nil
}

type Color int

const (
	None Color = iota
	BLACK
	WHITE
)

func (c Color) String() string {
	if c == BLACK {
		return "BLACK"
	}
	if c == WHITE {
		return "WHITE"
	}
	return "None"
}

func (c Color) Opposite() Color {
	if c == BLACK {
		return WHITE
	}
	if c == WHITE {
		return BLACK
	}
	return None
}

type Place struct {
	x int
	y int
}

func (p Place) next(d direction) Place {
	if d == UP {
		return Place{x: p.x, y: p.y - 1}
	}
	if d == DOWN {
		return Place{x: p.x, y: p.y + 1}
	}
	if d == LEFT {
		return Place{x: p.x - 1, y: p.y}
	}
	if d == RIGHT {
		return Place{x: p.x + 1, y: p.y}
	}
	if d == UP_LEFT {
		return Place{x: p.x - 1, y: p.y - 1}
	}
	if d == UP_RIGHT {
		return Place{x: p.x + 1, y: p.y - 1}
	}
	if d == DOWN_LEFT {
		return Place{x: p.x - 1, y: p.y + 1}
	}
	if d == DOWN_RIGHT {
		return Place{x: p.x + 1, y: p.y + 1}
	}
	panic("invalid direction")
}

type direction int

const (
	UP direction = iota
	DOWN
	LEFT
	RIGHT
	UP_RIGHT
	UP_LEFT
	DOWN_RIGHT
	DOWN_LEFT
)

type Command interface {
	WhoDoes() Color
	Commit(b *Board) error
	Describe() string
}
