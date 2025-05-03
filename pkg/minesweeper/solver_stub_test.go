//go:build !race

package minesweeper

import "testing"

func TestSolverConcurrencySafe(t *testing.T) {
	t.Skip("Test needs to be run with -race flag")
}
