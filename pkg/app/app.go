//go:generate fyne bundle --package assets --prefix Resource -o ../../assets/bundle_generated.go ../../img/mine.png
//go:generate fyne bundle --prefix Resource -o ../../assets/bundle_generated.go -append ../../img/flag.png
package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	fApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/heathcliff26/go-minesweeper/pkg/minesweeper"
)

var TEXT_COLOR = color.White

var DEFAULT_DIFFICULTY = minesweeper.Difficulties()[minesweeper.DifficultyIntermediate]

type App struct {
	app          fyne.App
	main         fyne.Window
	Version      Version
	grid         *MinesweeperGrid
	difficulties []*fyne.MenuItem
}

func New() *App {
	app := fApp.New()
	version := getVersion(app)
	main := app.NewWindow(version.Name)

	a := &App{
		app:     app,
		main:    main,
		Version: version,
		grid:    NewMinesweeperGrid(DEFAULT_DIFFICULTY),
	}

	a.makeMenu()
	a.main.SetMainMenu(a.makeMenu())

	a.setContent()

	a.main.SetFixedSize(true)
	a.main.Show()

	return a
}

func (a *App) Run() {
	a.app.Run()
}

func (a *App) makeMenu() *fyne.MainMenu {
	difficulties := minesweeper.Difficulties()
	diffItems := make([]*fyne.MenuItem, 0, len(difficulties))
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
	a.difficulties = diffItems
	diffMenu := fyne.NewMenu("Difficulties", diffItems...)

	about := fyne.NewMenuItem("About", func() {
		vInfo := dialog.NewCustom(a.Version.Name, "close", getVersionContent(a.Version), a.main)
		vInfo.Show()
	})
	helpMenu := fyne.NewMenu("Help", about)

	return fyne.NewMainMenu(diffMenu, helpMenu)
}

func (a *App) setContent() {
	a.main.SetContent(a.grid.GetCanvasObject())
	a.main.Resize(a.main.Content().MinSize())
}
