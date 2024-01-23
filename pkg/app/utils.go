package app

import (
	"image/color"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Version struct {
	Name, Version, Commit, Go string
}

// Extract the git commit
func getVersion(app fyne.App) Version {
	var commit string
	buildinfo, _ := debug.ReadBuildInfo()
	for _, item := range buildinfo.Settings {
		if item.Key == "vcs.revision" {
			commit = item.Value
			break
		}
	}
	if len(commit) > 7 {
		commit = commit[:7]
	}

	metadata := app.Metadata()

	name, _ := strings.CutSuffix(metadata.Name, ".exe")
	if name == "" {
		name = filepath.Base(os.Args[0])
	}

	return Version{
		Name:    name,
		Version: "v" + metadata.Version,
		Commit:  commit,
		Go:      runtime.Version(),
	}
}

// Create the content for the version dialog
func getVersionContent(v Version) fyne.CanvasObject {
	r1 := make([]fyne.CanvasObject, 3)
	r2 := make([]fyne.CanvasObject, 3)
	r1[0] = canvas.NewText("Version:", TEXT_COLOR)
	r2[0] = canvas.NewText(v.Version, TEXT_COLOR)
	r1[1] = canvas.NewText("Commit:", TEXT_COLOR)
	r2[1] = canvas.NewText(v.Commit, TEXT_COLOR)
	r1[2] = canvas.NewText("Go:", TEXT_COLOR)
	r2[2] = canvas.NewText(v.Go, TEXT_COLOR)

	row1 := container.NewVBox(r1...)
	row2 := container.NewVBox(r2...)

	return container.NewPadded(container.NewHBox(row1, row2))
}

// Create a line for the border
func makeBorderStrip() fyne.CanvasObject {
	rec := canvas.NewRectangle(color.White)
	rec.SetMinSize(fyne.NewSize(1, 1))
	return rec
}

// Wrap the objects in a box with border lines
func newBorder(content ...fyne.CanvasObject) fyne.CanvasObject {
	top := makeBorderStrip()
	left := makeBorderStrip()
	bottom := makeBorderStrip()
	right := makeBorderStrip()
	border := container.NewBorder(top, bottom, left, right, content...)
	return container.NewPadded(border)
}

func newGridLabel(text string) *canvas.Text {
	label := canvas.NewText(text, GridLabelColor)
	label.TextSize = GridLabelSize
	label.TextStyle.Bold = true
	return label
}
