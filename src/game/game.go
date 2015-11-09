package game

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"board"
)

type playerListener interface {
	getMove() board.Move
	player() board.Player
}

type ioListener struct {
	Reader *bufio.Reader
	Player board.Player
}

func (i ioListener) player() board.Player {
	return i.Player
}

func (i ioListener) getMove() board.Move {
	move := board.Move{}
	for {
		fmt.Printf("Player %s. Enter Move (ex: 1,1) :\n", i.Player)
		text, err := i.Reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		m_int, err2 := parseMove(text, i.Player)
		if err2 != nil {
			fmt.Println(err2)
			continue
		} else {
			move = m_int
			break
		}
	}
	return move
}

func parseMove(s string, f board.Player) (m board.Move, e error) {
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

func getInput(player <-chan playerListener, moves chan<- board.Move) {
	for {
		select {
		// Pop the current player out of the channel.
		case p := <-player:
			// Get her current move
			moves <- p.getMove()
		}
	}
}

func runBoard(moves <-chan board.Move, player chan<- playerListener, ended chan<- gameState) {
	cross := ioListener{Reader: bufio.NewReader(os.Stdin), Player: board.CROSS}
	circle := ioListener{Reader: bufio.NewReader(os.Stdin), Player: board.CIRCLE}
	// First player is cross, then circle.
	player <- cross
	i := 0
	nextPlayer := func() playerListener {
		i += 1
		if i%2 == 1 {
			return circle
		}
		return cross
	}
	b := board.Board{
		Finished: false,
	}
	for {
		select {
		case move := <-moves:
			err := b.UpdateBoard(move)
			if err != nil {
				// If an error occurred while updating the board, we skip back to the player issuing the move.
				fmt.Println("Invalid move, discarding next move to go back to player.")
				_ = nextPlayer()
				player <- nextPlayer()
				continue
			}
			b.PrintBoard()
			finished, winner := b.HasEnded()
			if finished {
				g := gameState{}
				g.finished = true
				// The winning player will be the one in the bottom of the channel, so we get twice on the channel.
				g.winner = winner
				ended <- g
				return
			}
			player <- nextPlayer()
		}
	}
}

func Main() {
	moves := make(chan board.Move)
	player := make(chan playerListener, 2)
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
