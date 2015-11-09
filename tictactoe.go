package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"board"
)

func getMove(s string, f board.Player) (m board.Move, e error) {
	text := strings.TrimSuffix(s, "\n")
	ss := strings.Split(text, ",")
	if len(ss) != 2 {
		e = errors.New("Error parsing move. Returning default move.")
		return
	}
	d, e1 := strconv.Atoi(ss[0])
	si, e2 := strconv.Atoi(ss[1])
	if e1 != nil || e2 != nil {
		fmt.Println(e1)
		fmt.Println(e2)
	}

	m = board.Move{Down: d, Side: si, Value: f}

	if !m.IsValidMove() {
		e = errors.New("Please provide a valid move.")
	}
	return
}

type gameState struct {
	finished bool
	winner   board.Player
}

func getInput(player <-chan board.Player, moves chan<- board.Move) {
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case f := <-player:
			fmt.Printf("Player %s. Enter Move (ex: 1,1) :\n", f)
			move := board.Move{}
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

func runBoard(moves chan board.Move, player chan board.Player, ended chan gameState) {
	f := board.CROSS
	player <- f
	b := board.Board{
		Finished: false,
	}
	for {
		select {
		case move := <-moves:
			err := b.UpdateBoard(move)
			if err != nil {
				fmt.Println(err)
				player <- f
				continue
			}
			b.PrintBoard()
			if b.HasEnded() {
				g := gameState{}
				g.finished = true
				g.winner = f
				ended <- g
			}
			f = f.SwitchPlayer()
			player <- f
		}
	}
}

func main() {
	moves := make(chan board.Move)
	player := make(chan board.Player)
	ended := make(chan gameState)
	go getInput(player, moves)
	go runBoard(moves, player, ended)
	select {
	case b := <-ended:
		if b.finished {
			fmt.Printf("Player %s won!", b.winner)
		}
	}

}
