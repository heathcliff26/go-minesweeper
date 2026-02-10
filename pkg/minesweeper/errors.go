package minesweeper

import "fmt"

type ErrDifficultyDimension struct {
	row int
	col int
}

func NewErrDifficultyDimension(row, col int) error {
	return &ErrDifficultyDimension{row, col}
}

func (e *ErrDifficultyDimension) Error() string {
	return fmt.Sprintf("Rows and Columns need to be between %d and %d, got %dx%d", DifficultyRowColMin, DifficultyRowColMax, e.row, e.col)
}

type ErrDifficultyMineCount struct {
	mines int
}

func NewErrDifficultyMineCount(mines int) error {
	return &ErrDifficultyMineCount{mines}
}

func (e *ErrDifficultyMineCount) Error() string {
	return fmt.Sprintf("The number of mines need to be between %d and %.1f %% of the total number of cells, got %d", DifficultyMineMin, DifficultyMineMaxPercentage*100, e.mines)
}

type ErrCreateUnsolvableGame struct {
	iterations int
}

func NewErrCreateUnsolvableGame(iterations int) error {
	return &ErrCreateUnsolvableGame{iterations}
}

func (e *ErrCreateUnsolvableGame) Error() string {
	return fmt.Sprintf("Could not create a solvable game within %d iterations", e.iterations)
}
