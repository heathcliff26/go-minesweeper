package minesweeper

import (
	"testing"

	"github.com/heathcliff26/go-minesweeper/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewGameWithSafePos(t *testing.T) {
	tMatrix := Difficulties()
	tMatrix = append(tMatrix, Difficulty{
		Name:  "InvertedExpert",
		Row:   30,
		Col:   16,
		Mines: 99,
	})
	p := NewPos(1, 1)

	for _, d := range tMatrix {
		t.Run(d.Name, func(t *testing.T) {
			g := NewGameWithSafePos(d, p)

			assert := assert.New(t)

			assert.Equal(d.Row, len(g.Field), "Should have the given number of rows")
			assert.Equal(d.Col, len(g.Field[0]), "Should have the given number of columns")
			assert.Equal(d, g.Difficulty, "Should have the given difficulty")
			assert.False(g.GameOver, "Should not be Game Over")
			assert.False(g.GameWon, "Game should not be won")
			assert.NotEqual(Mine, g.Field[p.X][p.Y].Content, "Safe position should not be a mine")

			mines := 0
			g.walkField(func(x, y int) {
				if g.Field[x][y].Content == Mine {
					mines++
				} else {
					assert.Equalf(g.countNearbyMines(NewPos(x, y)), int(g.Field[x][y].Content), "(%d, %d) should have the given number of neighboring mines", x, y)
				}
			})
			assert.Equal(d.Mines, mines, "Should have the given number of mines")
		})
	}
}

func TestNewGameWithSafeArea(t *testing.T) {
	tMatrix := Difficulties()
	tMatrix = append(tMatrix, Difficulty{
		Name:  "InvertedExpert",
		Row:   30,
		Col:   16,
		Mines: 99,
	})
	p := NewPos(1, 1)

	for _, d := range tMatrix {
		t.Run(d.Name, func(t *testing.T) {
			g := NewGameWithSafeArea(d, p)

			assert := assert.New(t)

			assert.Equal(d.Row, len(g.Field), "Should have the given number of rows")
			assert.Equal(d.Col, len(g.Field[0]), "Should have the given number of columns")
			assert.Equal(d, g.Difficulty, "Should have the given difficulty")
			assert.False(g.GameOver, "Should not be Game Over")
			assert.False(g.GameWon, "Game should not be won")
			for x := -1; x < 2; x++ {
				for y := -1; y < 2; y++ {
					assert.NotEqualf(Mine, g.Field[p.X+x][p.Y+y].Content, "Safe area should not be a mine, pos=(%d, %d)", p.X+x, p.Y+y)
				}
			}

			mines := 0
			g.walkField(func(x, y int) {
				if g.Field[x][y].Content == Mine {
					mines++
				} else {
					assert.Equalf(g.countNearbyMines(NewPos(x, y)), int(g.Field[x][y].Content), "(%d, %d) should have the given number of neighboring mines", x, y)
				}
			})
			assert.Equal(d.Mines, mines, "Should have the given number of mines")
		})
	}
}

func TestCheckField(t *testing.T) {
	minefield := [][]FieldContent{
		{Mine, 2, Mine, 1, 0, 0, 0, 0},
		{1, 2, 1, 2, 1, 2, 1, 1},
		{2, 2, 1, 1, Mine, 2, Mine, 1},
		{Mine, Mine, 1, 1, 1, 2, 1, 1},
		{Mine, 3, 1, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 2, 2},
		{0, 0, 0, 0, 0, 1, Mine, Mine},
	}
	minefield2 := [][]FieldContent{
		{Mine, 2, Mine, 1, 0, 0, 0, 0, 1, Mine},
		{1, 2, 1, 2, 1, 2, 1, 1, 2, 2},
		{2, 2, 1, 1, Mine, 2, Mine, 1, 1, Mine},
		{Mine, Mine, 1, 1, 1, 2, 1, 1, 1, 1},
		{Mine, 3, 1, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 2, 2, 1, 0},
		{0, 0, 0, 0, 0, 1, Mine, Mine, 1, 0},
	}
	minefield3 := [][]FieldContent{
		{Mine, 2, Mine, 1, 0, 0, 0, 0},
		{1, 2, 1, 2, 1, 2, 1, 1},
		{2, 2, 1, 1, Mine, 2, Mine, 1},
		{Mine, Mine, 1, 1, 1, 2, 1, 1},
		{Mine, 3, 1, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 2, 2},
		{0, 0, 0, 0, 0, 2, Mine, Mine},
		{0, 0, 0, 0, 0, 2, Mine, Mine},
		{0, 0, 0, 0, 0, 1, 3, Mine},
	}
	tMatrix := []struct {
		Name       string
		Pos        Pos
		GameOver   bool
		Difficulty Difficulty
		Minefield  [][]FieldContent
		Result     [][]bool
	}{
		{
			Name:       "Mine",
			Pos:        NewPos(0, 0),
			GameOver:   true,
			Difficulty: difficulties[DifficultyClassic],
			Minefield:  minefield,
		},
		{
			Name:       "Number",
			Pos:        NewPos(0, 1),
			Difficulty: difficulties[DifficultyClassic],
			Minefield:  minefield,
		},
		{
			Name:       "RevealMultipleFields1",
			Pos:        NewPos(0, 6),
			Minefield:  minefield,
			Difficulty: difficulties[DifficultyClassic],
			Result: [][]bool{
				{false, false, false, true, true, true, true, true},
				{false, false, false, true, true, true, true, true},
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
			},
		},
		{
			Name:       "RevealMultipleFields2",
			Pos:        NewPos(7, 3),
			Minefield:  minefield,
			Difficulty: difficulties[DifficultyClassic],
			Result: [][]bool{
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
				{false, false, true, true, true, true, true, true},
				{false, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, false, false},
			},
		},
		{
			Name:      "MoreColThanRow",
			Pos:       NewPos(4, 9),
			Minefield: minefield2,
			Difficulty: Difficulty{
				Name:  "MoreColThanRow",
				Row:   8,
				Col:   10,
				Mines: 11,
			},
			Result: [][]bool{
				{false, false, false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false, false, false},
				{false, false, true, true, true, true, true, true, true, true},
				{false, true, true, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, false, false, true, true},
			},
		},
		{
			Name:      "MoreRowThanCol",
			Pos:       NewPos(7, 3),
			Minefield: minefield3,
			Difficulty: Difficulty{
				Name:  "MoreRowThanCol",
				Row:   10,
				Col:   8,
				Mines: 12,
			},
			Result: [][]bool{
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
				{false, false, false, false, false, false, false, false},
				{false, false, true, true, true, true, true, true},
				{false, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, true, true},
				{true, true, true, true, true, true, false, false},
				{true, true, true, true, true, true, false, false},
				{true, true, true, true, true, true, false, false},
			},
		},
	}

	tMatrix[0].Result = utils.Make2D[bool](tMatrix[0].Difficulty.Row, tMatrix[0].Difficulty.Col)
	tMatrix[0].Result[0][0] = true

	tMatrix[1].Result = utils.Make2D[bool](tMatrix[1].Difficulty.Row, tMatrix[1].Difficulty.Col)
	tMatrix[1].Result[0][1] = true

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			d := tCase.Difficulty
			g := NewGameWithSafePos(d, tCase.Pos)

			g.walkField(func(x, y int) {
				g.Field[x][y].Content = tCase.Minefield[x][y]
			})
			g.CheckField(tCase.Pos)

			assert := assert.New(t)

			assert.Equal(tCase.GameOver, g.GameOver)

			g.walkField(func(x, y int) {
				assert.Equalf(tCase.Result[x][y], g.Field[x][y].Checked, "(%d, %d) Content: %d", x, y, g.Field[x][y].Content)
			})
		})
	}

	tMatrix2 := []struct {
		Name      string
		Win, Loss bool
	}{
		{"GameOver", false, true},
		{"GameWon", true, false},
	}

	for _, tCase := range tMatrix2 {
		t.Run(tCase.Name, func(t *testing.T) {
			d := difficulties[DifficultyIntermediate]
			g := NewGameWithSafePos(d, NewPos(0, 0))

			g.GameWon = tCase.Win
			g.GameOver = tCase.Loss

			s1 := g.Status()
			s2 := g.CheckField(NewPos(0, 0))

			assert.Equal(t, s1, s2, "Status should not change for a finished game")
		})
	}
}

func TestOutOfBounds(t *testing.T) {
	g := NewGameWithSafePos(difficulties[DifficultyExpert], NewPos(2, 2))
	d := g.Difficulty

	assert := assert.New(t)

	g.walkField(func(x, y int) {
		p := NewPos(x, y)
		assert.Falsef(g.OutOfBounds(p), "%v should be within the field", p)

		p.X = -1
		assert.Truef(g.OutOfBounds(p), "%v should be out of bounds", p)
		p.X = d.Row
		assert.Truef(g.OutOfBounds(p), "%v should be out of bounds", p)
		p.X = x

		p.Y = -1
		assert.Truef(g.OutOfBounds(p), "%v should be out of bounds", p)
		p.Y = d.Col
		assert.Truef(g.OutOfBounds(p), "%v should be out of bounds", p)
	})
}

func TestStatus(t *testing.T) {
	tMatrix := []struct {
		Name string
		Loss bool
		Win  bool
	}{
		{"GameInProgress", false, false},
		{"GameWon", false, true},
		{"GameOver", true, false},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			p := NewPos(1, 1)
			d := difficulties[DifficultyExpert]
			g := NewGameWithSafePos(d, p)

			g.GameOver = tCase.Loss
			g.GameWon = tCase.Win

			g.UpdateStatus()

			assert := assert.New(t)

			g.walkField(func(x, y int) {
				if tCase.Loss || tCase.Win {
					if !assert.Equal(g.Field[x][y], g.status.Field[x][y], "Fields should match") {
						t.FailNow()
					}
					return
				}

				if !assert.Equal(Unknown, g.status.Field[x][y].Content, "Field should be unknown") {
					t.FailNow()
				}
			})

			assert.Same(g.status, g.Status())
		})
	}
}

func TestDetectVictory(t *testing.T) {
	d := difficulties[DifficultyExpert]
	g := NewGameWithSafePos(d, NewPos(0, 0))

	assert := assert.New(t)

	assert.False(g.GameWon)
	assert.False(g.GameOver)

	g.walkField(func(x, y int) {
		g.Field[x][y].Checked = g.Field[x][y].Content != Mine
	})

	g.UpdateStatus()

	assert.True(g.GameWon)
	assert.True(g.status.GameWon)

	assert.False(g.GameOver)
	assert.False(g.status.GameOver)
}

func TestReplay(t *testing.T) {
	p := NewPos(1, 1)
	g := NewGameWithSafePos(difficulties[DifficultyBeginner], p)

	g.CheckField(p)
	g.Replay()

	assert := assert.New(t)

	assert.True(g.IsReplay(), "Game should be a replay")
	assert.False(g.GameOver, "Should not be Game Over")
	assert.False(g.GameWon, "Game should not be won")
	assert.NotEqual(Mine, g.Field[p.X][p.Y].Content, "Safe position should not be a mine")
	g.walkField(func(x, y int) {
		assert.Falsef(g.Field[x][y].Checked, "(%d, %d) should not be checked", x, y)
	})
}
