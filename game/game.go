package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Cell int

const (
	X Cell = iota
	O
	None
)

func (c Cell) String() string {
	var ch rune
	switch c {
	case X:
		ch = 'X'
	case O:
		ch = 'O'
	case None:
		ch = '-'
	}
	return fmt.Sprintf("%c", ch)
}

type TicTacToe struct {
	Cells [][]Cell
	WonBy Cell
}

type Position struct {
	Row int
	Col int
}

func ParsePosition(posStr string) (Position, error) {
	rowCol := strings.Split(posStr, ",")
	if len(rowCol) != 2 {
		return Position{}, fmt.Errorf("Invalid format, must have a comma: %v", posStr)
	}
	row, err := strconv.ParseInt(strings.TrimSpace(rowCol[0]), 10, 32)
	if err != nil {
		return Position{}, fmt.Errorf("%v is not a valid number", rowCol[0])
	}
	col, err := strconv.ParseInt(strings.TrimSpace(rowCol[1]), 10, 32)
	if err != nil {
		return Position{}, fmt.Errorf("%v is not a valid number", rowCol[1])
	}
	return Position{int(row), int(col)}, nil
}

func (p Position) Shift(vector Position) Position {
	return Position{
		p.Row + vector.Row,
		p.Col + vector.Col,
	}
}

func NewTicTacToe() *TicTacToe {
	return &TicTacToe{
		Cells: [][]Cell{
			{None, None, None},
			{None, None, None},
			{None, None, None},
		},
		WonBy: None,
	}
}

func (t *TicTacToe) CellAt(p Position) (Cell, error) {
	if p.Row > len(t.Cells) || p.Col > len(t.Cells[p.Row]) {
		return None, fmt.Errorf("Position out of bounds " + fmt.Sprintf("%+v", p))
	}
	return t.Cells[p.Row][p.Col], nil
}

func (t *TicTacToe) Play(ch Cell, p Position) (Cell, error) {
	cell, err := t.CellAt(p)
	if err != nil {
		return None, err
	}
	if cell != None {
		return None, fmt.Errorf("Cell already played")
	}
	t.Cells[p.Row][p.Col] = ch
	return t.DetectWin(), nil
}

// Accepts an ascii representation of the game in the format:
//
//         XXO
//         OO-
//         OXX
//
// "-" is None
//
// In the above game, we except DetectWin to be true for O
func ParseGame(gameStr string) (*TicTacToe, error) {
	game := NewTicTacToe()
	cells := strings.Split(strings.TrimSpace(gameStr), "\n")
	for row := 0; row < 3; row++ {
		rowStr := strings.TrimSpace(cells[row])
		for col := 0; col < 3; col++ {
			switch rowStr[col] {
			case 'X':
				game.Cells[row][col] = X
			case 'O':
				game.Cells[row][col] = O
			case '-':
				game.Cells[row][col] = None
			default:
				return nil, fmt.Errorf("Invalid character in string: %v", rowStr[col])
			}
		}
	}
	return game, nil
}

func (t *TicTacToe) String() string {
	retRunes := []rune{}
	for row := 0; row < 3; row++ {
		if row >= 1 {
			retRunes = append(retRunes, '\n')
		}
		for col := 0; col < 3; col++ {
			var ch rune
			switch t.Cells[row][col] {
			case X:
				ch = 'X'
			case O:
				ch = 'O'
			case None:
				ch = '-'
			default:
				panic("Invalid Cell value")
			}
			retRunes = append(retRunes, ch)
		}
	}
	return string(retRunes)
}

func (t *TicTacToe) DetectWin() Cell {
	winConditions := [][]Position{
		{
			{0, 0}, {0, 1},
		},
		{
			{1, 0}, {0, 1},
		},
		{
			{2, 0}, {0, 1},
		},
		{
			{0, 0}, {1, 0},
		},
		{
			{0, 1}, {1, 0},
		},
		{
			{0, 2}, {1, 0},
		},
		{
			{0, 0}, {1, 1},
		},
		{
			{2, 0}, {-1, 1},
		},
	}
	possibleWinner := None
	for _, condition := range winConditions {
		cell, _ := t.CellAt(condition[0])
		if cell == None {
			continue
		}
		possibleWinner = cell
		secondPosition := condition[0].Shift(condition[1])
		if cell, _ := t.CellAt(secondPosition); cell != possibleWinner {
			continue
		}
		thirdPosition := secondPosition.Shift(condition[1])
		if cell, _ := t.CellAt(thirdPosition); cell != possibleWinner {
			continue
		}
		return possibleWinner
	}
	return None
}

type Ultimate struct {
	innerGames [][]*TicTacToe
	metaGame   *TicTacToe
}

func NewUltimate() *Ultimate {
	return &Ultimate{
		innerGames: [][]*TicTacToe{
			{NewTicTacToe(), NewTicTacToe(), NewTicTacToe()},
			{NewTicTacToe(), NewTicTacToe(), NewTicTacToe()},
			{NewTicTacToe(), NewTicTacToe(), NewTicTacToe()},
		},
		metaGame: NewTicTacToe(),
	}
}

func (u *Ultimate) CellAt(p Position) (Cell, error) {
	return u.metaGame.CellAt(p)
}

func (u *Ultimate) GameAt(p Position) *TicTacToe {
	return u.innerGames[p.Row][p.Col]
}

func (u *Ultimate) String(oFormatter func(...interface{}) string,
	xFormatter func(...interface{}) string,
	normalFormatter func(...interface{}) string) string {
	ret := `
+-----++-----++-----+
| %s%s%s || %s%s%s || %s%s%s |
| %s%s%s || %s%s%s || %s%s%s |
| %s%s%s || %s%s%s || %s%s%s |
+-----++-----++-----+
| %s%s%s || %s%s%s || %s%s%s |
| %s%s%s || %s%s%s || %s%s%s |
| %s%s%s || %s%s%s || %s%s%s |
+-----++-----++-----+
| %s%s%s || %s%s%s || %s%s%s |
| %s%s%s || %s%s%s || %s%s%s |
| %s%s%s || %s%s%s || %s%s%s |
+-----++-----++-----+
`
	toFmt := []interface{}{}
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			gamePosition := Position{row / 3, col / 3}
			cellPosition := Position{row % 3, col % 3}
			game := u.GameAt(gamePosition)
			cell, _ := game.CellAt(cellPosition)
			win := game.DetectWin()
			if win == O {
				toFmt = append(toFmt, oFormatter(cell.String()))
			} else if win == X {
				toFmt = append(toFmt, xFormatter(cell.String()))
			} else {
				toFmt = append(toFmt, normalFormatter(cell.String()))
			}
		}
	}
	return fmt.Sprintf(ret, toFmt...)
}

func (u *Ultimate) Play(outer Position, inner Position, cell Cell) (Cell, error) {
	metaCell, err := u.metaGame.CellAt(outer)
	if err != nil {
		return None, err
	}
	if metaCell != None {
		return None, fmt.Errorf("Outer game %+v already won by %v", outer, metaCell.String())
	}
	winner, err := u.GameAt(outer).Play(cell, inner)
	if err != nil {
		return None, err
	}
	if winner != None {
		metaWinner, err := u.metaGame.Play(cell, outer)
		if err != nil {
			panic("Precondition failed- should not have been able to play inner game: " + err.Error())
		}
		if metaWinner != None {
			return metaWinner, nil
		}
	}
	return None, nil
}

type UltimateRunner struct {
	Game         *Ultimate
	CurrentTurn  Cell
	GamePosition *Position
}

func NewRunner() *UltimateRunner {
	return &UltimateRunner{
		Game:         NewUltimate(),
		CurrentTurn:  O,
		GamePosition: nil,
	}
}

func (u *UltimateRunner) Run(quit <-chan struct{}) {
	reader := bufio.NewReader(os.Stdin)

	gameColor := color.New(color.FgBlue)

	playerOColor := color.New(color.FgGreen, color.Bold)
	playerXColor := color.New(color.FgMagenta, color.Bold)

	errColor := color.New(color.FgRed, color.Bold)
	for {
		fmt.Print(u.Game.String(playerOColor.SprintFunc(), playerXColor.SprintFunc(), gameColor.SprintFunc()))

		var playerColor *color.Color
		if u.CurrentTurn == O {
			playerColor = playerOColor
		} else {
			playerColor = playerXColor
		}

		fmt.Printf("It's %s's turn\n", playerColor.SprintfFunc()("%s", u.CurrentTurn.String()))

		if u.GamePosition == nil {
			fmt.Printf("You may play anywhere! Enter the game coordinates as X,Y you'd like to play\n")
			in, err := reader.ReadString('\n')
			if err != nil {
				errColor.Printf("%v\n", err)
				continue
			}
			pos, err := ParsePosition(in)
			if err != nil {
				errColor.Printf("%v", err)
				continue
			}
			u.GamePosition = &pos
		}

		fmt.Printf("Enter move for game at %d,%d as X,Y coordinates\n", u.GamePosition.Row, u.GamePosition.Col)

		in, err := reader.ReadString('\n')
		if err != nil {
			errColor.Printf("%v\n", err)
			continue
		}
		pos, err := ParsePosition(in)
		if err != nil {
			errColor.Printf("%v", err)
			continue
		}

		cell, err := u.Game.Play(*u.GamePosition, pos, u.CurrentTurn)
		if err != nil {
			errColor.Printf("%v\n", err)
			continue
		}
		if cell != None {
			playerColor.Printf("%c Wins!", u.CurrentTurn)
			os.Exit(0)
			return
		}
		// see if the current player pointed the next player at a completed game
		if u.Game.GameAt(pos).DetectWin() != None {
			u.GamePosition = nil
		} else {
			u.GamePosition = &pos
		}

		// switch the current player
		if u.CurrentTurn == O {
			u.CurrentTurn = X
		} else {
			u.CurrentTurn = O
		}
	}
}
