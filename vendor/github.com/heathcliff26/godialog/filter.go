package godialog

type FileFilter struct {
	// A short description of this filter.
	// Example: "Images (*.png, *.jpg)"
	Description string
	// A list of file extensions to filter.
	// Example: []string{".png", ".jpg"}
	Extensions []string
}

type FileFilters []FileFilter

// Returns a list of all file extensions from all filters.
func (ff FileFilters) Extensions() []string {
	var extensions []string
	for _, filter := range ff {
		extensions = append(extensions, filter.Extensions...)
	}
	return extensions
}
