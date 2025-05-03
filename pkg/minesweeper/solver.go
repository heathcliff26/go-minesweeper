package minesweeper

import (
	"log/slog"
	"slices"
	"sync"
)

type Solver struct {
	game Game

	mines     []Pos
	nextSteps []Pos

	lock sync.RWMutex
}

func NewSolver(g Game) *Solver {
	return &Solver{
		game: g,
	}
}

func (s *Solver) NextSteps() []Pos {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.nextSteps
}

func (s *Solver) KnownMines() []Pos {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.mines
}

func (s *Solver) Update() {
	s.lock.Lock()
	defer s.lock.Unlock()

	status := s.game.Status()
	if status == nil || status.GameOver() || status.GameWon() {
		return
	}

	if s.mines == nil {
		s.mines = make([]Pos, 0, s.game.Difficulty().Mines)
	}
	nextSteps := make([]Pos, 0, 25)
	for _, p := range s.nextSteps {
		if !status.Field[p.X][p.Y].Checked {
			nextSteps = append(nextSteps, p)
		}
	}
	s.nextSteps = nextSteps

	i := 0
	oldLenMines := -1
	oldLenSafe := -1
	for len(s.mines) > oldLenMines || len(s.nextSteps) > oldLenSafe {
		i++
		oldLenMines = len(s.mines)
		oldLenSafe = len(s.nextSteps)
		walkField(func(x, y int) {
			if !status.Field[x][y].Checked || status.Field[x][y].Content <= 0 {
				return
			}

			unchecked := FieldContent(0)
			mines := FieldContent(0)
			safePos := FieldContent(0)
			newPos := make([]Pos, 0, 8)
			for m := -1; m < 2; m++ {
				for n := -1; n < 2; n++ {
					p := NewPos(x+m, y+n)
					if OutOfBounds(p, s.game.Difficulty()) {
						continue
					}
					if !status.Field[p.X][p.Y].Checked {
						unchecked++
						if slices.Contains(s.mines, p) {
							mines++
						} else if !slices.Contains(s.nextSteps, p) {
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
			if unchecked-safePos == status.Field[x][y].Content && mines != status.Field[x][y].Content {
				slog.Debug("Assisted Mode: Found mines near field",
					slog.String("pos", NewPos(x, y).String()),
					slog.Int("mines", int(mines)),
					slog.Int("unchecked", int(unchecked)),
					slog.Int("content", int(status.Field[x][y].Content)),
					slog.Any("newPos", newPos),
				)
				s.mines = append(s.mines, newPos...)
			} else if mines == status.Field[x][y].Content && unchecked > status.Field[x][y].Content {
				slog.Debug("Assisted Mode: Found safe positions near field",
					slog.String("pos", NewPos(x, y).String()),
					slog.Int("mines", int(mines)),
					slog.Int("unchecked", int(unchecked)),
					slog.Int("content", int(status.Field[x][y].Content)),
					slog.Any("newPos", newPos),
				)
				s.nextSteps = append(s.nextSteps, newPos...)
			}
		}, s.game.Difficulty().Row, s.game.Difficulty().Col)
	}

	slog.Debug("Assisted Mode: Mines and safe Positions found", slog.Any("mines", s.mines), slog.Any("safe", s.nextSteps), slog.Int("iterations", i))
}

// Autosolve the game.
// Returns if the game can be autosolved
func (s *Solver) Autosolve(startPos Pos) bool {
	s.game.CheckField(startPos)

	s.Update()

	for safePos := s.NextSteps(); len(safePos) > 0 && !s.game.Won(); safePos = s.NextSteps() {
		for _, p := range safePos {
			s.game.CheckField(p)
		}
		s.Update()
	}

	return s.game.Won()
}
