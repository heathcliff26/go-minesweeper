package minesweeper

const (
	DifficultyClassic      = iota
	DifficultyBeginner     = iota
	DifficultyIntermediate = iota
	DifficultyExpert       = iota
)

type Difficulty struct {
	Name  string
	Size  GridSize
	Mines int
}

type GridSize struct {
	Row, Col int
}

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

func Difficulties() []Difficulty {
	list := make([]Difficulty, len(difficulties))
	copy(list, difficulties)
	return list
}
