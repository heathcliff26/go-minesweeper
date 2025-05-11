package app

import (
	"log/slog"
	"os"
	"time"

	"fyne.io/fyne/v2"
	fApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/heathcliff26/go-minesweeper/pkg/app/locations"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
	"github.com/heathcliff26/godialog"
	fallbackfyne "github.com/heathcliff26/godialog/fallback/fyne"
)

const DEFAULT_DIFFICULTY = minesweeper.DifficultyBeginner

const DEFAULT_GAME_ALGORITHM = GameAlgorithmSafeArea

const DEFAULT_AUTOSOLVE_DELAY = 500 * time.Millisecond

// Used to change the new app function for testing
var newApp = fApp.New

var saveFileFilters = godialog.FileFilters{
	{
		Description: "Save File (*" + minesweeper.SaveFileExtension + ")",
		Extensions:  []string{minesweeper.SaveFileExtension},
	},
}

// Struct representing the current app.
// There should only ever be a single instance during runtime.
type App struct {
	app            fyne.App
	main           fyne.Window
	Version        Version
	grid           *MinesweeperGrid
	difficulties   []*fyne.MenuItem
	gameMenu       []*fyne.MenuItem
	assistedMode   *fyne.MenuItem
	gameAlgorithms []*fyne.MenuItem
	filedialog     godialog.FileDialog
}

// Create a new App
func New() *App {
	preferences, err := LoadPreferences()
	if err != nil {
		slog.Info("Failed to load preferences, falling back to defaults", "err", err)
	} else {
		slog.Debug("Loaded preferences", "preferences", preferences)
	}

	app := newApp()
	version := getVersion(app)
	main := app.NewWindow(version.Name)
	app.Settings().SetTheme(mainTheme{})

	fd := godialog.NewFileDialog()
	fd.SetFilters(saveFileFilters)
	fd.SetFallback(fallbackfyne.NewFyneFallbackDialog(app))

	a := &App{
		app:        app,
		main:       main,
		Version:    version,
		filedialog: fd,
	}
	a.main.SetTitle(version.Name)
	a.makeMenu(preferences)
	a.NewGrid(preferences.Difficulty())
	a.setGameAlgorithm(preferences.GameAlgorithm)

	a.main.SetFixedSize(true)
	a.main.Show()

	saveDir, err := locations.SaveFolder()
	if err != nil {
		dialog.ShowError(err, main)
	}
	a.filedialog.SetInitialDirectory(saveDir)

	return a
}

// Simply calls app.Run()
func (a *App) Run() {
	a.app.Run()

	preferences := CreatePreferencesFromApp(a)
	err := preferences.Save()
	if err != nil {
		slog.Error("Failed to save preferences for next time", "err", err)
	} else {
		slog.Info("Saved preferences")
	}
}

// Create the main menu bar
func (a *App) makeMenu(preferences Preferences) {
	// Can't assign grid functions directly, as the instance of grid may change
	newGameOption := fyne.NewMenuItem("New", func() {
		a.grid.NewGame()
	})
	replayOption := fyne.NewMenuItem("Replay", func() {
		a.grid.Replay()
	})
	loadOption := fyne.NewMenuItem("Load", a.loadSave)
	saveOption := fyne.NewMenuItem("Save", a.saveGame)
	a.gameMenu = []*fyne.MenuItem{newGameOption, replayOption, fyne.NewMenuItemSeparator(), loadOption, saveOption}
	gameMenu := fyne.NewMenu("Game", a.gameMenu...)

	difficulties := minesweeper.Difficulties()
	diffItems := make([]*fyne.MenuItem, 0, len(difficulties)+2)
	for _, d := range difficulties {
		item := fyne.NewMenuItem(d.Name, nil)
		item.Action = func() {
			if item.Checked {
				return
			}
			for _, i := range a.difficulties {
				i.Checked = (i.Label == d.Name)
			}
			a.NewGrid(d)
		}
		item.Checked = (d == preferences.Difficulty())
		diffItems = append(diffItems, item)
	}
	diffItems = append(diffItems, fyne.NewMenuItemSeparator())
	diffItems = append(diffItems, fyne.NewMenuItem("Custom", a.customDifficultyDialog))
	a.difficulties = diffItems
	diffMenu := fyne.NewMenu("Difficulties", diffItems...)

	a.assistedMode = fyne.NewMenuItem("      Assisted Mode", func() {
		a.assistedMode.Checked = !a.assistedMode.Checked
		a.grid.AssistedMode = a.assistedMode.Checked
		if a.grid.AssistedMode && a.grid.Game != nil {
			go a.grid.updateFromStatus(a.grid.Game.Status())
		}
	})
	a.assistedMode.Checked = preferences.AssistedMode
	a.gameAlgorithms = make([]*fyne.MenuItem, 3)
	a.gameAlgorithms[0] = fyne.NewMenuItem("Safe Position", func() {
		a.setGameAlgorithm(GameAlgorithmSafePos)
	})
	a.gameAlgorithms[1] = fyne.NewMenuItem("Safe Area", func() {
		a.setGameAlgorithm(GameAlgorithmSafeArea)
	})
	a.gameAlgorithms[2] = fyne.NewMenuItem("Solvable", func() {
		a.setGameAlgorithm(GameAlgorithmSolvable)
	})
	gameAlgorithmSubMenu := fyne.NewMenuItem("Creation Algorithm", nil)
	gameAlgorithmSubMenu.ChildMenu = fyne.NewMenu("Creation Algorithm", a.gameAlgorithms...)
	autosolve := fyne.NewMenuItem("Autosolve", func() {
		go func() {
			if !a.grid.Autosolve(DEFAULT_AUTOSOLVE_DELAY) {
				dialog.ShowInformation("Autosolve", "Failed to run autosolve, please ensure that a game is currently running.", a.main)
			}
		}()
	})
	optionsMenu := fyne.NewMenu("Options", a.assistedMode, gameAlgorithmSubMenu, autosolve)

	hint := fyne.NewMenuItem("Hint", func() {
		go func() {
			if !a.grid.Hint() {
				dialog.NewInformation("No hint found", "Could not find any hints to give.", a.main).Show()
			}
		}()
	})
	about := fyne.NewMenuItem("About", func() {
		vInfo := dialog.NewCustom(a.Version.Name, "close", getVersionContent(a.Version), a.main)
		vInfo.Show()
	})
	helpMenu := fyne.NewMenu("Help", hint, about)

	a.main.SetMainMenu(fyne.NewMainMenu(gameMenu, diffMenu, optionsMenu, helpMenu))
}

// Update the content of the app and resize the window to make it fit
func (a *App) setContent() {
	content := container.NewPadded(a.grid.GetCanvasObject())
	content.Resize(content.MinSize())

	a.main.SetContent(content)
	a.main.Resize(content.MinSize())
}

// Show a dialog for setting a custom difficulty
func (a *App) customDifficultyDialog() {
	mines := minesweeper.DifficultyMineMin
	row, col := minesweeper.DifficultyRowColMin, minesweeper.DifficultyRowColMin

	mineItem := widget.NewFormItem("Mines", widget.NewEntryWithData(binding.IntToString(binding.BindInt(&mines))))
	rowItem := widget.NewFormItem("Rows", widget.NewEntryWithData(binding.IntToString(binding.BindInt(&row))))
	colItem := widget.NewFormItem("Columns", widget.NewEntryWithData(binding.IntToString(binding.BindInt(&col))))

	content := widget.NewForm(mineItem, rowItem, colItem)

	diffDialog := dialog.NewCustomConfirm("Custom Difficulty", "ok", "cancel", content, func(ok bool) {
		if !ok {
			return
		}
		d, err := minesweeper.NewCustomDifficulty(mines, row, col)
		if err != nil {
			dialog.ShowError(err, a.main)
			return
		}

		for _, i := range a.difficulties {
			i.Checked = (i.Label == "Custom")
		}
		a.NewGrid(d)
	}, a.main)
	diffDialog.Show()
}

func (a *App) loadSave() {
	a.filedialog.Open("Open Savegame", a.loadSaveCallback)
}

func (a *App) loadSaveCallback(path string, err error) {
	if err != nil {
		dialog.ShowError(err, a.main)
		return
	}
	if path == "" {
		return
	}

	save, err := minesweeper.LoadSave(path)
	if err != nil {
		dialog.ShowError(err, a.main)
		return
	}

	for _, i := range a.difficulties {
		i.Checked = (i.Label == save.Data.Difficulty.Name)
	}
	fyne.DoAndWait(func() {
		a.NewGrid(save.Data.Difficulty)
	})

	a.grid.Game = save.Game()
}

func (a *App) saveGame() {
	if a.grid.Game == nil {
		d := dialog.NewInformation("Can't save game", "You need to first start a game before you can save it.", a.main)
		d.Show()
		return
	}
	a.filedialog.Save("Save Game", a.saveGameCallback)
}

func (a *App) saveGameCallback(path string, err error) {
	if err != nil {
		dialog.ShowError(err, a.main)
		return
	}
	if path == "" {
		return
	}

	// Ensure that empty files created by fynes file dialog are removed before creating the save.
	// This ensures that no empty files without ".sav" extension are created.
	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		dialog.ShowError(err, a.main)
		return
	}

	save, err := a.grid.Game.ToSave()
	if err != nil {
		dialog.ShowError(err, a.main)
		return
	}

	err = save.Save(path)
	if err != nil {
		dialog.ShowError(err, a.main)
		return
	}
}

func (a *App) NewGrid(d minesweeper.Difficulty) {
	a.grid = NewMinesweeperGrid(d, a.assistedMode.Checked)
	for i, item := range a.gameAlgorithms {
		if item.Checked {
			a.grid.GameAlgorithm = i
		}
	}
	a.setContent()
}

func (a *App) setGameAlgorithm(id int) {
	for i, item := range a.gameAlgorithms {
		item.Checked = i == id
	}
	a.grid.GameAlgorithm = id
}
