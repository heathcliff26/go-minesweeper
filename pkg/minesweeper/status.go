package minesweeper

// Status contains the current state of the game known to the player.
// As such it will always be a copy and needs to be it's own type, despite the
// overlapping similarities.
// It does not support any of the functions that Game does.
// It is safe to write to Status, as it is merely a copy.
type Status struct {
	Field    [][]Field
	gameOver bool
	gameWon  bool
}

// Returns if the game is lost
func (s *Status) GameOver() bool {
	return s.gameOver
}

// Returns if the game is won
func (s *Status) GameWon() bool {
	return s.gameWon
}
