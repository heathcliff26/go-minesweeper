package fyne

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/heathcliff26/godialog"
)

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

type GenericURICloser interface {
	Close() error
	URI() fyne.URI
}

func callCallback(cb godialog.DialogCallback, uri GenericURICloser, err error) {
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
