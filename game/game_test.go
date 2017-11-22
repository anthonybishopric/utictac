package game

import (
	"fmt"
	"testing"
)

func TestParseGame(t *testing.T) {
	gameState := `
XO-
-O-
XO-
`
	toe, err := ParseGame(gameState)
	if err != nil {
		t.Fatalf(err.Error())
	}
	str := toe.String()
	if str != `XO-
-O-
XO-` {
		t.Fatalf("Incorrect string for game state: %s", str)
	}
}

func TestDetectWin(t *testing.T) {
	tests := []struct {
		state  string
		winner Cell
	}{
		{`
    X-X
    OOO
    -X-
    `,
			O,
		},
		{
			`
    XOX
    XOO
    X-O
    `,
			X,
		},
		{
			`
    OXO
    XXX
    OO-
    `,
			X,
		},
		{
			`
    -XO
    O-X
    XOX
    `,
			None,
		},
		{
			`
      OXO
      XOX
      O-X
      `,
			O,
		},
		{
			`
      XO-
      -XO
      O-X
      `,
			X,
		},
	}
	for _, test := range tests {
		game, err := ParseGame(test.state)
		if err != nil {
			t.Errorf("Test %+v failed to parseGame: %v", test, err)
		}
		if game.DetectWin() != test.winner {
			t.Errorf("Test %+v had unexpected winner %v", test, game.DetectWin())
		}
	}
}

func TestUltimateString(t *testing.T) {
	innerGame1, err := ParseGame(`XXO
O-O
X-O`)
	innerGame2, err := ParseGame(`X-O
-XO
XOX`)
	if err != nil {
		t.Fatalf("Error during test setup: %v", err)
	}
	ult := NewUltimate()
	ult.innerGames[1][2] = innerGame1
	ult.innerGames[2][2] = innerGame2
	expected := `
+-----++-----++-----+
| --- || --- || --- |
| --- || --- || --- |
| --- || --- || --- |
+-----++-----++-----+
| --- || --- || XXO |
| --- || --- || O-O |
| --- || --- || X-O |
+-----++-----++-----+
| --- || --- || X-O |
| --- || --- || -XO |
| --- || --- || XOX |
+-----++-----++-----+
`

	if expected != ult.String(fmt.Sprint, fmt.Sprint, fmt.Sprint) {
		t.Fatalf("Got \n%s\n instead", ult.String(fmt.Sprint, fmt.Sprint, fmt.Sprint))
	}
}
