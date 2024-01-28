package game

import (
	"errors"
)

const BOARDSIZE int = 8
const FIRST Color = BLACK

type Board struct {
	history []Command
	pieces  map[Place]Color
}

func NewBoard() Board {
	var b Board
	b.pieces = make(map[Place]Color)
	for y := 0; y < BOARDSIZE; y++ {
		for x := 0; x < BOARDSIZE; x++ {
			b.placePiece(NewPlace(x, y), None)
		}
	}
	// 初期配置
	b.placePiece(NewPlace(3, 3), WHITE)
	b.placePiece(NewPlace(4, 3), BLACK)
	b.placePiece(NewPlace(3, 4), BLACK)
	b.placePiece(NewPlace(4, 4), WHITE)

	return b
}

// ディープコピーする。
func (b Board) Clone() Board {
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
		logs = append(logs, c.Describe())
	}
	return logs
}

func (b Board) Pieces() map[Place]Color {
	return b.Clone().pieces
}

func (b Board) WhoesTurn() Color {
	if len(b.history) == 0 {
		return FIRST
	}
	return b.history[len(b.history)-1].WhoDoes().Opposite()
}

func (b Board) TurnCount() int {
	return len(b.history)
}

// 盤面をnum回巻き戻す
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

// Undoと同様。
// 失敗時はpanicする。
func (b *Board) MustUndo(num int) {
	err := b.Undo(num)
	if err != nil {
		panic(err)
	}
}

// 得点を取得する。
// Board.Score()[BLACK] = 黒の得点
// Board.Score()[WHITE] = 白の得点
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

// 盤面に対してCommand.Commitを実行する。
// Command.Commitが失敗した場合は、盤面に作用しない。
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

// MakeEffectと同様。
// Command.Commitが失敗した場合はpanicする。
func (b *Board) MustMakeEffect(c Command) {
	err := b.MakeEffect(c)
	if err != nil {
		panic(err)
	}
}

func (b Board) Hints(c Color) []PlacePiece {
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

	return ps
}

func (b Board) scan(expect Color, d direction, p path) *path {
	currentColor, err := b.GetPiece(p.current)
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

func (p path) intoFlips() []flipPiece {
	var f []flipPiece
	for _, l := range p.log {
		f = append(f, flipPiece{place: l})
	}
	return f
}

func (b Board) IsInBound(p Place) bool {
	return p.x >= 0 || p.x < BOARDSIZE || p.y >= 0 || p.y < BOARDSIZE
}

func (b *Board) placePiece(place Place, c Color) error {
	if !b.IsInBound(place) {
		return errors.New("out of range")
	}
	b.pieces[place] = c
	return nil
}

func (b Board) GetPiece(place Place) (*Color, error) {
	if !b.IsInBound(place) {
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

func NewPlace(x, y int) Place {
	return Place{x, y}
}

func (p Place) next(d direction) Place {
	directions := map[direction]Place{
		UP:         {x: p.x, y: p.y - 1},
		DOWN:       {x: p.x, y: p.y + 1},
		LEFT:       {x: p.x - 1, y: p.y},
		RIGHT:      {x: p.x + 1, y: p.y},
		UP_LEFT:    {x: p.x - 1, y: p.y - 1},
		UP_RIGHT:   {x: p.x + 1, y: p.y - 1},
		DOWN_LEFT:  {x: p.x - 1, y: p.y + 1},
		DOWN_RIGHT: {x: p.x + 1, y: p.y + 1},
	}
	return directions[d]
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
