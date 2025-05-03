//go:build race

package minesweeper

import (
	"sync"
	"testing"
)

func TestSolverConcurrencySafe(t *testing.T) {
	p := NewPos(0, 0)
	g := NewGameWithSafePos(Difficulties()[0], p)
	g.CheckField(p)

	s := NewSolver(g)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		for i := 0; i < 1000; i++ {
			_ = s.NextSteps()
		}
	}()
	go func() {
		defer wg.Done()

		for i := 0; i < 1000; i++ {
			_ = s.KnownMines()
		}
	}()
	go func() {
		defer wg.Done()

		for i := 0; i < 1000; i++ {
			s.Update()
		}
	}()

	wg.Wait()
}
