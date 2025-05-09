package filedialog

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestSetDialogLocationToDir(t *testing.T) {
	assert := assert.New(t)

	a := test.NewApp()
	w := a.NewWindow("Test")
	d := dialog.NewFileOpen(nil, w)

	// Test with a valid directory
	dir := "."
	assert.NoError(setDialogLocationToDir(dir, d), "Should set the location for a valid directory")

	// Test with an no directory
	dir = ""
	assert.Error(setDialogLocationToDir(dir, d), "Should return an error for a no directory")

	// Test with non-existing directory
	dir = "/this/path/does/not/exist"
	assert.Error(setDialogLocationToDir(dir, d), "Should return an error for a non-existing directory")
}

type MockGenericURICloser struct {
	uri         fyne.URI
	calledClose bool
}

func (m *MockGenericURICloser) Close() error {
	m.calledClose = true
	return nil
}

func (m *MockGenericURICloser) URI() fyne.URI {
	return m.uri
}

func TestCallCallback(t *testing.T) {
	tMatrix := []struct {
		Name string
		Path string
		Err  error
	}{
		{
			Name: "Error",
			Err:  fmt.Errorf("error"),
		},
		{
			Name: "NoPath",
		},
		{
			Name: "ValidPath",
			Path: ".",
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			called := false

			cb := func(path string, err error) {
				called = true

				assert.Equal(tCase.Path, path, "Path should be equal")
				assert.Equal(tCase.Err, err, "Error should be equal")
			}

			if tCase.Path != "" {
				mUri := &MockGenericURICloser{}
				mUri.uri, _ = storage.ParseURI("file://" + tCase.Path)
				callCallback(cb, mUri, tCase.Err)
				assert.True(mUri.calledClose, "Close function of URI should be called")
			} else {
				callCallback(cb, nil, tCase.Err)
			}

			assert.True(called, "Callback should be called")
		})
	}
}
