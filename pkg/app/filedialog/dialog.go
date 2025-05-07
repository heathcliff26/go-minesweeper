package filedialog

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

const (
	DialogHeight = 800
	DialogWidth  = 600
)

type FileFilter map[string][]string

type GenericURICloser interface {
	Close() error
	URI() fyne.URI
}

// Use internal fyne file dialog to open a file.
func internalFileOpen(name string, startLocation string, filters FileFilter, cb func(string, error)) {
	w := fyne.CurrentApp().NewWindow(name)
	d := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
		// Ensure this runs in a goroutine as we call fyne.DoAndWait in the callback
		go callCallback(cb, uri, err)
	}, w)

	err := showFileDialog(startLocation, convertToExtensions(filters), d, w)
	if err != nil {
		cb("", err)
	}
}

// Use internal fyne file dialog to save a file.
func internalFileSave(name string, startLocation string, filters FileFilter, cb func(string, error)) {
	w := fyne.CurrentApp().NewWindow(name)
	d := dialog.NewFileSave(func(uri fyne.URIWriteCloser, err error) {
		// Ensure this runs in a goroutine as we call fyne.DoAndWait in the callback
		go callCallback(cb, uri, err)
	}, w)

	err := showFileDialog(startLocation, convertToExtensions(filters), d, w)
	if err != nil {
		cb("", err)
	}
}

// Set a file dialogs location to the given directory.
// When dir is empty, uses current directory.
// Returns error on failure.
func setDialogLocationToDir(dir string, d *dialog.FileDialog) error {
	uri, err := storage.ParseURI("file://" + dir)
	if err != nil {
		return fmt.Errorf("failed to parse URI: %w", err)
	}
	listURI, err := storage.ListerForURI(uri)
	if err != nil {
		return fmt.Errorf("failed to create lister for URI: %w", err)
	}
	d.SetLocation(listURI)

	return nil
}

func showFileDialog(startLocation string, extensions []string, d *dialog.FileDialog, w fyne.Window) error {
	d.SetFilter(storage.NewExtensionFileFilter(extensions))

	err := setDialogLocationToDir(startLocation, d)
	if err != nil {
		return err
	}

	d.SetOnClosed(func() {
		w.Close()
	})

	w.Resize(fyne.NewSize(DialogHeight, DialogWidth))
	w.SetFixedSize(true)
	d.Resize(fyne.NewSize(DialogHeight, DialogWidth))
	fyne.Do(func() {
		w.Show()
		d.Show()
	})

	return nil
}

func convertToExtensions(filters FileFilter) []string {
	var extensions []string
	for _, filter := range filters {
		extensions = append(extensions, filter...)
	}
	return extensions
}

func callCallback(cb func(string, error), uri GenericURICloser, err error) {
	if err != nil {
		cb("", err)
		return
	}
	if uri == nil {
		cb("", nil)
		return
	}
	defer uri.Close()

	path := uri.URI().Path()
	cb(path, nil)
}
