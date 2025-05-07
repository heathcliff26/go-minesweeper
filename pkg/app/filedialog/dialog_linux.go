//go:build linux

package filedialog

import "log/slog"

// Show a file open dialog in a new window and return path.
func FileOpen(name string, startLocation string, filters FileFilter, cb func(string, error)) {
	err := dbusFileChooser(FileChooserOpenFile, name, startLocation, filters)
	if err != nil {
		slog.Info("Failed to open os native file dialog, falling back to internal implementation", "error", err)
		internalFileOpen(name, startLocation, filters, cb)
		return
	}

	go cb(dbusWaitForResponse())
}

// Show a file save dialog in a new window and return path.
func FileSave(name string, startLocation string, filters FileFilter, cb func(string, error)) {
	err := dbusFileChooser(FileChooserSaveFile, name, startLocation, filters)
	if err != nil {
		slog.Info("Failed to open os native file dialog, falling back to internal implementation", "error", err)
		internalFileSave(name, startLocation, filters, cb)
		return
	}

	go cb(dbusWaitForResponse())
}
