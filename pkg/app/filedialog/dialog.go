package filedialog

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

const (
	DialogHeight = 800
	DialogWidth  = 600
)

type FileDialogResult struct {
	Path  string
	Error error
}

// Show a file open dialog in a new window and return path
func FileOpen(name string, startLocation string, extensions []string) (string, error) {
	var called bool
	w := fyne.CurrentApp().NewWindow(name)
	ch := make(chan FileDialogResult)
	d := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
		called = true
		if err != nil {
			ch <- FileDialogResult{"", err}
			return
		}
		if uri == nil {
			ch <- FileDialogResult{"", err}
			return
		}

		ch <- FileDialogResult{uri.URI().Path(), nil}
	}, w)

	err := showFileDialog(startLocation, extensions, d, w)
	if err != nil {
		return "", err
	}

	w.SetOnClosed(func() {
		if called {
			return
		}
		ch <- FileDialogResult{}
	})

	result := <-ch
	return result.Path, result.Error
}

// Show a file save dialog in a new window and return path
func FileSave(name string, startLocation string, extensions []string) (string, error) {
	var called bool
	w := fyne.CurrentApp().NewWindow(name)
	ch := make(chan FileDialogResult)
	d := dialog.NewFileSave(func(uri fyne.URIWriteCloser, err error) {
		called = true
		if err != nil {
			ch <- FileDialogResult{"", err}
			return
		}
		if uri == nil {
			ch <- FileDialogResult{"", err}
			return
		}

		ch <- FileDialogResult{uri.URI().Path(), nil}
		uri.Close()
	}, w)

	err := showFileDialog(startLocation, extensions, d, w)
	if err != nil {
		return "", err
	}

	w.SetOnClosed(func() {
		if called {
			return
		}
		ch <- FileDialogResult{}
	})

	result := <-ch
	return result.Path, result.Error
}

// Set a file dialogs location to the given directory.
// When dir is empty, uses current directory.
// Returns error on failure.
func setDialogLocationToDir(dir string, d *dialog.FileDialog) error {
	uri, err := storage.ParseURI("file://" + dir)
	if err != nil {
		return err
	}
	listURI, err := storage.ListerForURI(uri)
	if err != nil {
		return err
	}
	d.SetLocation(listURI)

	return nil
}

func showFileDialog(startLocation string, extensions []string, d *dialog.FileDialog, w fyne.Window) error {
	d.SetFilter(storage.NewExtensionFileFilter(extensions))

	if startLocation == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		startLocation = pwd
	}
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
	w.Show()
	d.Show()

	return nil
}
