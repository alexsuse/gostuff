package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Player int

const (
	EMPTY  Player = 0
	CROSS  Player = 1
	CIRCLE Player = 2
)

type Move struct {
	down  int
	side  int
	value Player
}

type Board struct {
	board    [3][3]Player
	finished bool
}

func isValidMove(m Move) bool {
	return m.down >= 0 && m.down <= 2 && m.side >= 0 && m.side <= 2 &&
		(m.value == CROSS || m.value == CIRCLE)
}

func getMove(s string, f Player) (m Move, e error) {
	text := strings.TrimSuffix(s, "\n")
	ss := strings.Split(text, ",")
	if len(ss) != 2 {
		e = errors.New("Error parsing move. Returning default move.")
		return
	}
	down, e1 := strconv.Atoi(ss[0])
	if e1 != nil {
		fmt.Println(e1)
	}
	side, e2 := strconv.Atoi(ss[1])
	if e2 != nil {
		fmt.Println(e2)
	}

	m.down = down
	m.side = side
	m.value = f

	if !isValidMove(m) {
		e = errors.New("Please provide a valid move.")
	}
	return
}

func updateBoard(b *Board, m Move) (e error) {
	if !isValidMove(m) {
		e = errors.New("Please provide a valid move.")
		return
	} else if b.board[m.down][m.side] != EMPTY {
		e = errors.New("Field already occupied.")
		return
	}
	b.board[m.down][m.side] = m.value
	return
}

func hasEnded(b Board) bool {
	return false
}

func makeLineString(s [3]Player) string {
	return fmt.Sprintf("%s %s %s", getPlayerString(s[0]), getPlayerString(s[1]),
		getPlayerString(s[2]))
}

func getPlayerString(f Player) string {
	if f == CROSS {
		return "x"
	}
	if f == CIRCLE {
		return "o"
	}
	return "."
}

func printBoard(b Board) {
	fmt.Println("")
	for i := 0; i <= 2; i++ {
		fmt.Println(makeLineString(b.board[i]))
	}
	fmt.Println("")
}

func switchPlayer(f Player) (r Player) {
	if f == CROSS {
		return CIRCLE
	} else if f == CIRCLE {
		return CROSS
	}
	return EMPTY
}

type gameState struct {
	finished bool
	winner   Player
}

func getInput(player chan Player, moves chan Move) {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case f := <-player:
			fmt.Printf("Player %s. Enter Move (ex: 1,1) :\n", getPlayerString(f))
			move := Move{}
			for {
				text, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					continue
				}
				m_int, err2 := getMove(text, f)
				if err2 != nil {
					fmt.Println(err2)
					continue
				} else {
					move = m_int
					break
				}
			}
			moves <- move
		}
	}
}

func runBoard(moves chan Move, player chan Player, ended chan gameState) {
	f := CROSS
	player <- f
	b := Board{
		finished: false,
	}
	for {
		select {
		case move := <-moves:
			err := updateBoard(&b, move)
			if err != nil {
				fmt.Println(err)
				player <- f
				continue
			}
			printBoard(b)
			if hasEnded(b) {
				g := gameState{}
				g.finished = true
				g.winner = f
				ended <- g
			}
			f = switchPlayer(f)
			player <- f
		}
	}
}

func main() {
	moves := make(chan Move)
	player := make(chan Player)
	ended := make(chan gameState)
	go getInput(player, moves)
	go runBoard(moves, player, ended)
	select {
	case b := <-ended:
		if b.finished {
			return
		} else {
			return
		}
	}

}
