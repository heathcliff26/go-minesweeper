//go:build windows

package filedialog

// Show a file open dialog in a new window and return path.
func FileOpen(name string, startLocation string, filters FileFilter, cb func(string, error)) {
	internalFileOpen(name, startLocation, filters, cb)
}

// Show a file save dialog in a new window and return path.
func FileSave(name string, startLocation string, filters FileFilter, cb func(string, error)) {
	internalFileSave(name, startLocation, filters, cb)
}
