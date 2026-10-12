//go:build android

package godialog

import (
	"fmt"
)

// Show a file open dialog in a new window and return path.
func (fd *fileDialog) Open(title string, cb DialogCallback) {
	go fd.open(title, cb)
}

// Show a file save dialog in a new window and return path.
func (fd *fileDialog) Save(title string, cb DialogCallback) {
	go fd.save(title, cb)
}

// The actual implementation for Open. Should be run in a goroutine.
func (fd *fileDialog) open(title string, cb DialogCallback) {
	if fd.fallback != nil {
		fd.fallback.Open(title, fd.InitialDirectory(), fd.filters, cb)
	} else {
		cb("", fmt.Errorf("android native file dialog is not supported and no fallback is set"))
	}
}

// The actual implementation for Save. Should be run in a goroutine.
func (fd *fileDialog) save(title string, cb DialogCallback) {
	if fd.fallback != nil {
		fd.fallback.Save(title, fd.InitialDirectory(), fd.filters, cb)
	} else {
		cb("", fmt.Errorf("android native file dialog is not supported and no fallback is set"))
	}
}
