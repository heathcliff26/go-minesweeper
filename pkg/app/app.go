//go:generate fyne bundle --package assets --prefix Resource -o ../../assets/bundle_generated.go ../../img/mine.png
//go:generate fyne bundle --prefix Resource -o ../../assets/bundle_generated.go -append ../../img/flag.png
package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	fApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
)

var TEXT_COLOR = color.White

var DEFAULT_DIFFICULTY = minesweeper.Difficulties()[minesweeper.DifficultyIntermediate]

// Used to change the new app function for testing
var newApp = fApp.New

// Struct representing the current app.
// There should only ever be a single instance during runtime.
type App struct {
	app          fyne.App
	main         fyne.Window
	Version      Version
	grid         *MinesweeperGrid
	difficulties []*fyne.MenuItem
}

// Create a new App
func New() *App {
	app := newApp()
	version := getVersion(app)
	main := app.NewWindow(version.Name)

	a := &App{
		app:     app,
		main:    main,
		Version: version,
		grid:    NewMinesweeperGrid(DEFAULT_DIFFICULTY),
	}

	a.main.SetTitle(version.Name)
	a.makeMenu()
	a.main.SetMainMenu(a.makeMenu())

	a.setContent()

	a.main.SetFixedSize(true)
	a.main.Show()

	return a
}

// Simply calls app.Run()
func (a *App) Run() {
	a.app.Run()
}

// Create the main menu bar
func (a *App) makeMenu() *fyne.MainMenu {
	// A quit item will be added by fyne automatically to this menu
	appMenu := fyne.NewMenu("App", fyne.NewMenuItemSeparator())
	appMenu.Items = appMenu.Items[1:]

	difficulties := minesweeper.Difficulties()
	diffItems := make([]*fyne.MenuItem, 0, len(difficulties)+2)
	for _, d := range difficulties {
		diff := d
		item := fyne.NewMenuItem(d.Name, nil)
		item.Action = func() {
			if item.Checked {
				return
			}
			for _, i := range a.difficulties {
				i.Checked = (i.Label == diff.Name)
			}
			a.grid = NewMinesweeperGrid(diff)
			a.setContent()
		}
		item.Checked = (d == DEFAULT_DIFFICULTY)
		diffItems = append(diffItems, item)
	}
	diffItems = append(diffItems, fyne.NewMenuItemSeparator())
	diffItems = append(diffItems, fyne.NewMenuItem("Custom", a.customDifficultyDialog))
	a.difficulties = diffItems
	diffMenu := fyne.NewMenu("Difficulties", diffItems...)

	about := fyne.NewMenuItem("About", func() {
		vInfo := dialog.NewCustom(a.Version.Name, "close", getVersionContent(a.Version), a.main)
		vInfo.Show()
	})
	helpMenu := fyne.NewMenu("Help", about)

	return fyne.NewMainMenu(appMenu, diffMenu, helpMenu)
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

	mineLabel := canvas.NewText("Mines", TEXT_COLOR)
	mineEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&mines)))
	rowLabel := canvas.NewText("Rows", TEXT_COLOR)
	rowEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&row)))
	colLabel := canvas.NewText("Columns", TEXT_COLOR)
	colEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&col)))

	content := container.NewGridWithColumns(2, mineLabel, mineEntry, rowLabel, rowEntry, colLabel, colEntry)
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
		a.grid = NewMinesweeperGrid(d)
		a.setContent()
	}, a.main)
	diffDialog.Show()
}
