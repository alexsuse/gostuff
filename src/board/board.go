package board

import (
	"errors"
	"fmt"
)

type Player int

const (
	EMPTY  Player = 0
	CROSS  Player = 1
	CIRCLE Player = 2
)

type Move struct {
	Down  int
	Side  int
	Value Player
}

type Board struct {
	board    [3][3]Player
	Finished bool
}

func (m Move) IsValidMove() bool {
	return m.Down >= 0 && m.Down <= 2 && m.Side >= 0 && m.Side <= 2 &&
		(m.Value == CROSS || m.Value == CIRCLE)
}

func (b *Board) UpdateBoard(m Move) (e error) {
	if !m.IsValidMove() {
		e = errors.New("Please provide a valid move.")
		return
	} else if b.board[m.Down][m.Side] != EMPTY {
		e = errors.New("Field already occupied.")
		return
	}
	b.board[m.Down][m.Side] = m.Value
	return
}

func (b Board) checkLines() (bool, Player) {
	for i := 0; i < 3; i++ {
		if b.board[i][0] == b.board[i][1] && b.board[i][1] == b.board[i][2] {
			if b.board[i][0] != EMPTY {
				return true, b.board[i][0]
			}
		}
	}
	return false, EMPTY
}

func (b Board) checkColumns() (bool, Player) {
	for i := 0; i < 3; i++ {
		if b.board[0][i] == b.board[1][i] && b.board[1][i] == b.board[2][i] {
			if b.board[0][i] != EMPTY {
				return true, b.board[0][i]
			}
		}
	}
	return false, EMPTY
}

func (b Board) checkDiagonals() (bool, Player) {
	if b.board[0][0] == b.board[1][1] && b.board[1][1] == b.board[2][2] ||
		b.board[2][0] == b.board[1][1] && b.board[1][1] == b.board[0][2] {
		if b.board[1][1] != EMPTY {
			return true, b.board[1][1]
		}
	}
	return false, EMPTY
}

func (b Board) HasEnded() (bool, Player) {
	ended, winner := b.checkLines()
	if ended {
		return ended, winner
	}
	ended, winner = b.checkColumns()
	if ended {
		return ended, winner
	}
	ended, winner = b.checkDiagonals()
	if ended {
		return ended, winner
	}
	return false, EMPTY
}

func (p Player) String() string {
	if p == CROSS {
		return "x"
	} else if p == CIRCLE {
		return "o"
	}
	return "."
}

func makeLineString(p []Player) string {
	var s string
	for _, p := range p {
		s += fmt.Sprintf(" %s", p)
	}
	return s + "\n"
}

func (f Player) SwitchPlayer() (r Player) {
	if f == CROSS {
		return CIRCLE
	} else if f == CIRCLE {
		return CROSS
	}
	return EMPTY
}

func (b Board) PrintBoard() {
	var s string
	for _, pl := range b.board[:] {
		s += makeLineString(pl[:])
	}
	fmt.Printf("%s", s)
}
