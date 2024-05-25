package minesweeper

import (
	"log/slog"
	"slices"
	"sync"
)

// Status contains the current state of the game known to the player.
// As such it will always be a copy and needs to be it's own type, despite the
// overlapping similarities.
// It does not support any of the functions that Game does.
// It is safe to write to Status, as it is merely a copy.
type Status struct {
	Field      [][]Field
	gameOver   bool
	gameWon    bool
	difficulty Difficulty

	actionsUpdated bool
	actions        Actions
	mutex          sync.Mutex

	// Used only for unit-tests
	updateActionsCalled func()
}

type Actions struct {
	Mines   []Pos
	SafePos []Pos
}

// Returns if the game is lost
func (s *Status) GameOver() bool {
	return s.gameOver
}

// Returns if the game is won
func (s *Status) GameWon() bool {
	return s.gameWon
}

// Returns the position of all obvious mines
func (s *Status) ObviousMines() []Pos {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.actionsUpdated {
		s.updateActions()
	}
	return s.actions.Mines
}

// Returns the position of all obvious safe positions
func (s *Status) ObviousSafePos() []Pos {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.actionsUpdated {
		s.updateActions()
	}
	return s.actions.SafePos
}

// Calculate the next actions based on the current status
func (s *Status) updateActions() {
	if s.updateActionsCalled != nil {
		s.updateActionsCalled()
	}

	s.actionsUpdated = true
	if len(s.Field) == 0 || s.GameOver() || s.GameWon() {
		return
	}
	if s.actions.Mines == nil {
		s.actions.Mines = make([]Pos, 0, s.difficulty.Mines)
	}
	safePos := make([]Pos, 0, 25)
	for _, p := range s.actions.SafePos {
		if !s.Field[p.X][p.Y].Checked {
			safePos = append(safePos, p)
		}
	}
	s.actions.SafePos = safePos

	i := 0
	oldLenMines := -1
	oldLenSafe := -1
	for len(s.actions.Mines) > oldLenMines || len(s.actions.SafePos) > oldLenSafe {
		i++
		oldLenMines = len(s.actions.Mines)
		oldLenSafe = len(s.actions.SafePos)
		walkField(func(x, y int) {
			if !s.Field[x][y].Checked || s.Field[x][y].Content <= 0 {
				return
			}

			unchecked := FieldContent(0)
			mines := FieldContent(0)
			safePos := FieldContent(0)
			newPos := make([]Pos, 0, 8)
			for m := -1; m < 2; m++ {
				for n := -1; n < 2; n++ {
					p := NewPos(x+m, y+n)
					if OutOfBounds(p, s.difficulty) {
						continue
					}
					if !s.Field[p.X][p.Y].Checked {
						unchecked++
						if slices.Contains(s.actions.Mines, p) {
							mines++
						} else if !slices.Contains(s.actions.SafePos, p) {
							newPos = append(newPos, p)
						} else {
							safePos++
						}
					}
				}
			}
			if len(newPos) == 0 {
				return
			}
			if unchecked-safePos == s.Field[x][y].Content && mines != s.Field[x][y].Content {
				slog.Debug("Assisted Mode: Found mines near field",
					slog.String("pos", NewPos(x, y).String()),
					slog.Int("mines", int(mines)),
					slog.Int("unchecked", int(unchecked)),
					slog.Int("content", int(s.Field[x][y].Content)),
					slog.Any("newPos", newPos),
				)
				s.actions.Mines = append(s.actions.Mines, newPos...)
			} else if mines == s.Field[x][y].Content && unchecked > s.Field[x][y].Content {
				slog.Debug("Assisted Mode: Found safe positions near field",
					slog.String("pos", NewPos(x, y).String()),
					slog.Int("mines", int(mines)),
					slog.Int("unchecked", int(unchecked)),
					slog.Int("content", int(s.Field[x][y].Content)),
					slog.Any("newPos", newPos),
				)
				s.actions.SafePos = append(s.actions.SafePos, newPos...)
			}
		}, s.difficulty.Row, s.difficulty.Col)
	}

	slog.Debug("Assisted Mode: Mines and safe Positions found", slog.Any("mines", s.actions.Mines), slog.Any("safe", s.actions.SafePos), slog.Int("iterations", i))
}
