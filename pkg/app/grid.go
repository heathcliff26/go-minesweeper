package app

import (
	"image/color"
	"log/slog"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/heathcliff26/go-minesweeper/pkg/utils"
)

var (
	GridLabelColor = color.RGBA{240, 10, 20, alpha}
)

const (
	GridLabelSize float32 = 40

	ResetDefaultText          = "🙂"
	ResetGameOverText         = "☠"
	ResetGameWonText          = "😎"
	ResetTextSize     float32 = 40
)

const (
	GameAlgorithmSafePos = iota
	GameAlgorithmSafeArea
	GameAlgorithmSolvable
)

const (
	ChunkSize = 10
)

// Graphical display for a minesweeper game
type MinesweeperGrid struct {
	Tiles         [][]*Tile
	Difficulty    minesweeper.Difficulty
	Game          minesweeper.Game
	AssistedMode  bool
	GameAlgorithm int

	solver *minesweeper.Solver

	Timer       *Timer
	MineCount   *Counter
	ResetButton *Button

	lUpdate    sync.Mutex
	lGame      sync.Mutex
	lAutosolve sync.Mutex

	autosolveBreak chan bool
	autosolveDone  chan bool

	testChannel chan string
}

// Create a new grid suitable for the give difficulty
func NewMinesweeperGrid(d minesweeper.Difficulty, assistedMode bool) *MinesweeperGrid {
	tiles := utils.Make2D[*Tile](d.Row, d.Col)
	grid := &MinesweeperGrid{
		Tiles:         tiles,
		Difficulty:    d,
		AssistedMode:  assistedMode,
		GameAlgorithm: DEFAULT_GAME_ALGORITHM,
		Timer:         NewTimer(),
		MineCount:     NewCounter(d.Mines),
	}
	grid.ResetButton = NewButton(ResetDefaultText, color.RGBA{}, grid.NewGame)
	grid.ResetButton.Label.TextSize = ResetTextSize

	for x := 0; x < grid.Row(); x++ {
		for y := 0; y < grid.Col(); y++ {
			grid.Tiles[x][y] = NewTile(x, y, grid)
		}
	}

	return grid
}

// Get the graphical representation of the grid
func (g *MinesweeperGrid) GetCanvasObject() fyne.CanvasObject {
	mineCount := container.NewHBox(layout.NewSpacer(), container.NewCenter(newBorder(g.MineCount.Label)))
	reset := container.NewCenter(g.ResetButton)
	timer := container.NewHBox(container.NewCenter(newBorder(g.Timer.Label)), layout.NewSpacer())

	head := newBorder(container.NewGridWithColumns(3, mineCount, reset, timer))

	rows := make([]fyne.CanvasObject, len(g.Tiles))

	for x := 0; x < g.Row(); x++ {
		col := make([]fyne.CanvasObject, g.Col())
		for y := 0; y < g.Col(); y++ {
			col[y] = g.Tiles[x][y]
		}
		rows[x] = container.NewGridWithColumns(g.Col(), col...)
	}
	body := newBorder(container.NewGridWithRows(g.Row(), rows...))
	return container.NewVBox(head, body)
}

// Called by the child tiles to signal they have been tapped.
// Checks the given tile and then updates the display according to the new state.
// Starts a new game when no game is currently running.
func (g *MinesweeperGrid) TappedTile(pos minesweeper.Pos) {
	g.lGame.Lock()
	defer g.lGame.Unlock()

	if g.Game == nil {
		switch g.GameAlgorithm {
		case GameAlgorithmSafePos:
			g.Game = minesweeper.NewGameWithSafePos(g.Difficulty, pos)
		case GameAlgorithmSafeArea:
			g.Game = minesweeper.NewGameWithSafeArea(g.Difficulty, pos)
		case GameAlgorithmSolvable:
			g.Game = minesweeper.NewGameSolvable(g.Difficulty, pos)
		default:
			slog.Error("Unkown Algorithm for creating a new game", slog.Int("algorithm", g.GameAlgorithm))
			return
		}
	}
	if !g.Timer.Running() {
		g.Timer.Start()
	}

	slog.Info("Checking field", slog.String("pos", pos.String()))

	s, changed := g.Game.CheckField(pos)
	if changed {
		slog.Debug("Checked field, updating tiles")
		g.updateFromStatus(s)
	}

	if g.testChannel != nil {
		g.testChannel <- "TappedTile"
	}
}

// Called by the child tiles to reveal all neighbours when they have been double tapped
func (g *MinesweeperGrid) TapNeighbours(pos minesweeper.Pos) {
	g.lGame.Lock()
	defer g.lGame.Unlock()

	flags := minesweeper.FieldContent(0)
	posToCheck := make([]minesweeper.Pos, 0, 8)
	for m := -1; m < 2; m++ {
		for n := -1; n < 2; n++ {
			p := pos
			p.X += m
			p.Y += n
			if !g.OutOfBounds(p) {
				if g.Tiles[p.X][p.Y].Flagged() {
					flags++
					continue
				}
				if g.Tiles[p.X][p.Y].untappable() {
					continue
				}
				posToCheck = append(posToCheck, p)
			}
		}
	}

	if flags != g.Tiles[pos.X][pos.Y].Content() {
		return
	}

	for _, p := range posToCheck {
		g.Game.CheckField(p)
	}
	g.updateFromStatus(g.Game.Status())

	if g.testChannel != nil {
		g.testChannel <- "TapNeighbours"
	}
}

// Update the grid from the given status
func (g *MinesweeperGrid) updateFromStatus(s *minesweeper.Status) {
	if s == nil {
		return
	}

	g.lUpdate.Lock()
	defer g.lUpdate.Unlock()

	var wg sync.WaitGroup

	if s.GameOver() || s.GameWon() {
		switch {
		case s.GameWon():
			slog.Info("Win")
			g.ResetButton.SetText(ResetGameWonText)
		case s.GameOver():
			slog.Info("Game Over")
			g.ResetButton.SetText(ResetGameOverText)
		}
		g.Timer.Stop()
	} else if g.AssistedMode {
		if g.solver == nil {
			g.solver = minesweeper.NewSolver(g.Game)
		}
		g.solver.Update()

		wg.Add(2)
		slog.Debug("Creating Markers for Assisted Mode")
		go func() {
			defer wg.Done()
			for _, p := range g.solver.KnownMines() {
				g.Tiles[p.X][p.Y].Mark(HelpMarkingMine)
			}
		}()
		go func() {
			defer wg.Done()
			for _, p := range g.solver.NextSteps() {
				g.Tiles[p.X][p.Y].Mark(HelpMarkingSafe)
			}
		}()
	}

	for x := 0; x <= g.Row()/ChunkSize; x++ {
		chunkSizeX := ChunkSize
		if x == g.Row()/ChunkSize {
			if g.Row()%ChunkSize == 0 {
				continue
			} else {
				chunkSizeX = g.Row() % ChunkSize
			}
		}
		startX := x * ChunkSize

		for y := 0; y <= g.Col()/ChunkSize; y++ {
			chunkSizeY := ChunkSize
			if y == g.Col()/ChunkSize {
				if g.Col()%ChunkSize == 0 {
					continue
				} else {
					chunkSizeY = g.Col() % ChunkSize
				}
			}

			startY := y * ChunkSize

			wg.Add(1)
			go g.updateChunk(startX, startY, startX+chunkSizeX, startY+chunkSizeY, s, &wg)
		}
	}

	wg.Wait()
	slog.Debug("Finished Update")
}

// Update a single chunk defined by the given dimensions
func (g *MinesweeperGrid) updateChunk(startX, startY, endX, endY int, s *minesweeper.Status, wg *sync.WaitGroup) {
	defer wg.Done()
	for x := startX; x < endX; x++ {
		for y := startY; y < endY; y++ {
			if s.Field[x][y].Content == minesweeper.Unknown {
				continue
			}
			g.Tiles[x][y].SetField(s.Field[x][y])
		}
	}
}

// Return the number of rows in the grid
func (g *MinesweeperGrid) Row() int {
	return g.Difficulty.Row
}

// Return the number of columns in the grid
func (g *MinesweeperGrid) Col() int {
	return g.Difficulty.Col
}

// Start a new game
func (g *MinesweeperGrid) NewGame() {
	// Ensure that the game lock is released before calling reset.
	// This ensures that no deadlock occurs when autosolve is running.
	func() {
		g.lGame.Lock()
		defer g.lGame.Unlock()

		slog.Info("Preparing for new game")
		g.Game = nil
		g.solver = nil
	}()
	g.Reset()
}

// Replay the current game
func (g *MinesweeperGrid) Replay() {
	// Ensure that the game lock is released before calling reset.
	// This ensures that no deadlock occurs when autosolve is running.
	func() {
		g.lGame.Lock()
		defer g.lGame.Unlock()

		slog.Info("Preparing for replay of current game")
		if g.Game != nil {
			g.Game.Replay()
		}
	}()

	g.Reset()
}

// Reset Grid
func (g *MinesweeperGrid) Reset() {
	g.lAutosolve.Lock()
	defer g.lAutosolve.Unlock()

	if g.autosolveBreak != nil && g.autosolveDone != nil {
		close(g.autosolveBreak)
		<-g.autosolveDone
	}

	g.lUpdate.Lock()
	defer g.lUpdate.Unlock()

	for x := 0; x < g.Row(); x++ {
		for y := 0; y < g.Col(); y++ {
			g.Tiles[x][y].Reset()
		}
	}
	g.MineCount.SetCount(g.Difficulty.Mines)
	g.Timer.Reset()
	g.ResetButton.SetText(ResetDefaultText)
	g.ResetButton.Refresh()
	slog.Debug("Reset grid")
}

// Check if the given position is out of bounds.
// Calls Game.OutOfBounds(Pos)
func (g *MinesweeperGrid) OutOfBounds(p minesweeper.Pos) bool {
	return g.Game.OutOfBounds(p)
}

// Display a single hint on the grid.
// Returns false if no hint could be displayed.
func (g *MinesweeperGrid) Hint() bool {
	if !g.gameRunning() {
		return false
	}

	g.updateSolver()

	g.lUpdate.Lock()
	defer g.lUpdate.Unlock()

	for _, mine := range g.solver.KnownMines() {
		tile := g.Tiles[mine.X][mine.Y]
		if tile.Flagged() {
			continue
		}
		tile.Mark(HelpMarkingMine)
		return true
	}
	safePos := g.solver.NextSteps()
	if len(safePos) > 0 {
		g.Tiles[safePos[0].X][safePos[0].Y].Mark(HelpMarkingSafe)
		return true
	}
	return false
}

// Autosolve the current running game, delay is the time between steps
// Returns false if it can't run autosolve, otherwise true.
// If there are no steps to be taken, still returns true.
func (g *MinesweeperGrid) Autosolve(delay time.Duration) bool {
	if !g.gameRunning() {
		slog.Info("Autosolve failed because no game is running")
		return false
	}

	s := g.gameStatus()

	if s == nil {
		slog.Info("Autosolve failed because game does not have a status yet")
		return false
	}

	autosolveBreak := make(chan bool, 1)
	autosolveDone := make(chan bool, 1)
	alreadyRunning := false
	func() {
		g.lAutosolve.Lock()
		defer g.lAutosolve.Unlock()

		if g.autosolveBreak != nil && g.autosolveDone != nil {
			alreadyRunning = true
		} else {
			g.autosolveBreak = autosolveBreak
			g.autosolveDone = autosolveDone
		}
	}()

	if alreadyRunning {
		slog.Debug("Autosolve already running, aborting")
		return false
	}

	defer func() {
		close(autosolveDone)

		g.lAutosolve.Lock()
		defer g.lAutosolve.Unlock()

		g.autosolveBreak = nil
		g.autosolveDone = nil
	}()

	g.updateSolver()

	oldAssistedModeStatus := g.AssistedMode
	g.AssistedMode = true
	defer func() {
		g.AssistedMode = oldAssistedModeStatus
	}()
	g.updateFromStatus(s)

	for i, safePos := 0, g.solver.NextSteps(); len(safePos) > 0 && !s.GameOver() && !s.GameWon(); safePos = g.solver.NextSteps() {
		mines := g.solver.KnownMines()

		slog.Debug("Autosolve: Checking safe positions", slog.Int("iteration", i))
		for _, p := range safePos {
			if s.GameOver() || s.GameWon() {
				break
			}

			if g.Tiles[p.X][p.Y].Checked() {
				continue
			}
			g.TappedTile(p)
			time.Sleep(delay)
		}

		slog.Info("Autosolve: Flagging mines", slog.Int("iteration", i))
		for _, mine := range mines {
			if s.GameOver() || s.GameWon() {
				break
			}
			tile := g.Tiles[mine.X][mine.Y]
			tile.Flag(true)
		}

		select {
		case <-autosolveBreak:
			slog.Debug("Autosolve interrupted")
			return true
		default:
		}

		g.updateSolver()
		i++
	}

	if s.GameOver() || s.GameWon() {
		slog.Debug("Autosolve finished")
		return true
	}

	slog.Info("Autosolve: Flagging mines a final time")
	for _, mine := range g.solver.KnownMines() {
		tile := g.Tiles[mine.X][mine.Y]
		tile.Flag(true)
	}

	slog.Debug("Autosolve finished")
	return true
}

// Check if a game is currently running
func (g *MinesweeperGrid) gameRunning() bool {
	g.lGame.Lock()
	defer g.lGame.Unlock()

	if g.Game == nil {
		return false
	}
	return g.Game != nil && !g.Game.Won() && !g.Game.Lost()
}

// Update the autosolver.
// Initializes the solver if it is nil.
// Expects the game to be != nil
func (g *MinesweeperGrid) updateSolver() {
	g.lGame.Lock()
	defer g.lGame.Unlock()

	if g.solver == nil {
		g.solver = minesweeper.NewSolver(g.Game)
	}
	g.solver.Update()
}

// Return the current games status.
// Can return nil.
func (g *MinesweeperGrid) gameStatus() *minesweeper.Status {
	g.lGame.Lock()
	defer g.lGame.Unlock()

	if g.Game == nil {
		return nil
	}
	return g.Game.Status()
}
