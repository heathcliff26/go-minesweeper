package minesweeper

const (
	DifficultyClassic      = iota
	DifficultyBeginner     = iota
	DifficultyIntermediate = iota
	DifficultyExpert       = iota
)

// Represent a difficulty setting for the Game
type Difficulty struct {
	Name     string
	Row, Col int
	Mines    int
}

// Pre-defined difficulties
var difficulties []Difficulty = []Difficulty{
	{
		Name:  "Classic",
		Row:   8,
		Col:   8,
		Mines: 9,
	},
	{
		Name:  "Beginner",
		Row:   9,
		Col:   9,
		Mines: 10,
	},
	{
		Name:  "Intermediate",
		Row:   16,
		Col:   16,
		Mines: 40,
	},
	{
		Name:  "Expert",
		Row:   16,
		Col:   30,
		Mines: 99,
	},
}

// Exposes pre-defined difficulties in a way that does not allow the original array to be modified
func Difficulties() []Difficulty {
	list := make([]Difficulty, len(difficulties))
	copy(list, difficulties)
	return list
}
