package minesweeper

const (
	DifficultyClassic      = iota
	DifficultyBeginner     = iota
	DifficultyIntermediate = iota
	DifficultyExpert       = iota
)

// Represent a difficulty setting for the Game
type Difficulty struct {
	Name  string
	Size  GridSize
	Mines int
}

type GridSize struct {
	Row, Col int
}

// Pre-defined difficulties
var difficulties []Difficulty = []Difficulty{
	{
		Name:  "Classic",
		Size:  GridSize{8, 8},
		Mines: 9,
	},
	{
		Name:  "Beginner",
		Size:  GridSize{9, 9},
		Mines: 10,
	},
	{
		Name:  "Intermediate",
		Size:  GridSize{16, 16},
		Mines: 40,
	},
	{
		Name:  "Expert",
		Size:  GridSize{16, 30},
		Mines: 99,
	},
}

// Exposes pre-defined difficulties in a way that does not allow the original array to be modified
func Difficulties() []Difficulty {
	list := make([]Difficulty, len(difficulties))
	copy(list, difficulties)
	return list
}
