package filedialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

const (
	DialogHeight = 800
	DialogWidth  = 600
)

// Show a file open dialog in a new window and return path.
func FileOpen(name string, startLocation string, extensions []string, cb func(string, error)) {
	w := fyne.CurrentApp().NewWindow(name)
	d := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
		if err != nil {
			cb("", err)
			return
		}
		if uri == nil {
			cb("", err)
			return
		}

		cb(uri.URI().Path(), nil)
	}, w)

	err := showFileDialog(startLocation, extensions, d, w)
	if err != nil {
		cb("", err)
	}
}

// Show a file save dialog in a new window and return path.
func FileSave(name string, startLocation string, extensions []string, cb func(string, error)) {
	w := fyne.CurrentApp().NewWindow(name)
	d := dialog.NewFileSave(func(uri fyne.URIWriteCloser, err error) {
		if err != nil {
			cb("", err)
			return
		}
		if uri == nil {
			cb("", err)
			return
		}
		defer uri.Close()
		cb(uri.URI().Path(), nil)
	}, w)

	err := showFileDialog(startLocation, extensions, d, w)
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
