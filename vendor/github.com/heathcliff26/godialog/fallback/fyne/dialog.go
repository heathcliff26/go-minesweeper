package fyne

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/heathcliff26/godialog"
)

const (
	DefaultDialogHeight float32 = 800
	DefaultDialogWidth  float32 = 600
)

// Ensure FyneFallbackDialog implements godialog.FallbackDialog
var _ godialog.FallbackDialog = &FyneFallbackDialog{}

// Opens a file dialog in a new fyne window for the given app.
type FyneFallbackDialog struct {
	App    fyne.App
	Height float32
	Width  float32
}

func NewFyneFallbackDialog(app fyne.App) *FyneFallbackDialog {
	return &FyneFallbackDialog{
		App:    app,
		Height: DefaultDialogHeight,
		Width:  DefaultDialogWidth,
	}
}

// Shows the open file dialog and calls the callback asynchronously.
func (f *FyneFallbackDialog) Open(title string, initialDirectory string, filters godialog.FileFilters, cb godialog.DialogCallback) {
	if f.App == nil {
		go cb("", fmt.Errorf("cannot open file dialog: fyne.App is nil"))
		return
	}

	w := f.App.NewWindow(title)
	d := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
		// Ensure this runs in a goroutine as we call fyne.DoAndWait in the callback
		go callCallback(cb, uri, err)
	}, w)

	err := f.showFileDialog(initialDirectory, filters, d, w)
	if err != nil {
		go cb("", err)
	}
}

// Shows the save file dialog and calls the callback asynchronously.
func (f *FyneFallbackDialog) Save(title string, initialDirectory string, filters godialog.FileFilters, cb godialog.DialogCallback) {
	if f.App == nil {
		go cb("", fmt.Errorf("cannot open file dialog: fyne.App is nil"))
		return
	}

	w := f.App.NewWindow(title)
	d := dialog.NewFileSave(func(uri fyne.URIWriteCloser, err error) {
		// Ensure this runs in a goroutine as we call fyne.DoAndWait in the callback
		go callCallback(cb, uri, err)
	}, w)

	err := f.showFileDialog(initialDirectory, filters, d, w)
	if err != nil {
		go cb("", err)
	}
}

func (f *FyneFallbackDialog) showFileDialog(initialDirectory string, filters godialog.FileFilters, d *dialog.FileDialog, w fyne.Window) error {
	d.SetFilter(storage.NewExtensionFileFilter(filters.Extensions()))

	err := setDialogLocationToDir(initialDirectory, d)
	if err != nil {
		return err
	}

	d.SetOnClosed(func() {
		w.Close()
	})

	w.Resize(fyne.NewSize(f.Height, f.Width))
	w.SetFixedSize(true)
	d.Resize(fyne.NewSize(f.Height, f.Width))
	fyne.Do(func() {
		w.Show()
		d.Show()
	})

	return nil
}
