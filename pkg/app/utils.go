package app

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NOTE: The $Format strings are replaced during 'git archive' thanks to the
// companion .gitattributes file containing 'export-subst' in this same
// directory.  See also https://git-scm.com/docs/gitattributes
var gitCommit string = "$Format:%H$" // sha1 from git, output of $(git rev-parse HEAD)

func init() {
	initGitCommit()
}

func initGitCommit() {
	if strings.HasPrefix(gitCommit, "$Format") {
		var commit string
		buildinfo, _ := debug.ReadBuildInfo()
		for _, item := range buildinfo.Settings {
			if item.Key == "vcs.revision" {
				commit = item.Value
				break
			}
		}
		gitCommit = commit
	}
}

// Struct for containing the current version of the app
type Version struct {
	Name, Version, Commit, Go string
}

// Extract the version information from app
func getVersion(app fyne.App) Version {
	commit := gitCommit
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
	data := [][]string{
		{"Version:", v.Version},
		{"Commit:", v.Commit},
		{"Go:", v.Go},
	}

	versionTable := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("                ")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		},
	)

	versionTable.ShowHeaderRow = false
	versionTable.ShowHeaderColumn = false
	versionTable.StickyRowCount = len(data) - 1
	versionTable.StickyColumnCount = len(data[0]) - 1
	versionTable.HideSeparators = true

	return versionTable
}

// Wrap the objects in a box with border lines
func newBorder(content ...fyne.CanvasObject) fyne.CanvasObject {
	contentContainer := container.NewThemeOverride(container.NewPadded(content...), mainTheme{})
	border := widget.NewCard("", "", contentContainer)

	return container.NewThemeOverride(border, borderTheme{})
}

// Create a new label used in the grid, with preset color, text size and text style
func newGridLabel(text string) *canvas.Text {
	label := canvas.NewText(text, GridLabelColor)
	label.TextSize = GridLabelSize
	label.TextStyle.Bold = true
	return label
}
