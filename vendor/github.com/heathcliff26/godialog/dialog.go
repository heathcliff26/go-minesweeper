package godialog

// Ensure fileDialog implements the FileDialog interface.
var _ FileDialog = &fileDialog{}

type DialogCallback func(string, error)

type FallbackDialog interface {
	// Shows the open file dialog and calls the callback asynchronously.
	Open(title string, initialDirectory string, filters FileFilters, cb DialogCallback)
	// Shows the save file dialog and calls the callback asynchronously.
	Save(title string, initialDirectory string, filters FileFilters, cb DialogCallback)
}

// OS native file dialog. Allows to define a fallback implementation in case it does not work.
// callbacks are always called asynchronously.
type FileDialog interface {
	// Return the current file filters (nil if no filters are set).
	Filters() FileFilters
	// Add a new filter to the list of filters.
	AddFilter(filter FileFilter)
	// Set the file filters.
	SetFilters(filters FileFilters)
	// The current fallback dialog.
	// Returns nil if no fallback is set.
	Fallback() FallbackDialog
	// Set the fallback dialog in case the native dialog does not work.
	SetFallback(fallback FallbackDialog)
	// Set the initial directory for the file dialog.
	SetInitialDirectory(dir string)
	// Get the initial directory for the file dialog.
	InitialDirectory() string

	// Show a file open dialog in a new window and return path.
	// Runs asynchronously.
	Open(title string, cb DialogCallback)
	// Show a file save dialog in a new window and return path.
	// Runs asynchronously.
	Save(title string, cb DialogCallback)
}

// OS native file dialog. Allows to define a fallback implementation in case it does not work.
// File dialogs are always opened asynchronously.
type fileDialog struct {
	// The directory that the file dialog should open in.
	initialDirectory string
	filters          FileFilters
	fallback         FallbackDialog
}

// Create a new file dialog.
func NewFileDialog() FileDialog {
	return &fileDialog{}
}

// Return the current file filters (nil if no filters are set).
func (fd *fileDialog) Filters() FileFilters {
	return fd.filters
}

// Add a new filter to the list of filters.
func (fd *fileDialog) AddFilter(filter FileFilter) {
	fd.filters = append(fd.filters, filter)
}

// Set the file filters.
func (fd *fileDialog) SetFilters(filters FileFilters) {
	fd.filters = filters
}

// The current fallback dialog.
// Returns nil if no fallback is set.
func (fd *fileDialog) Fallback() FallbackDialog {
	return fd.fallback
}

// Set the fallback dialog in case the native dialog does not work.
func (fd *fileDialog) SetFallback(fallback FallbackDialog) {
	fd.fallback = fallback
}

// Set the initial directory for the file dialog.
func (fd *fileDialog) SetInitialDirectory(dir string) {
	fd.initialDirectory = dir
}

// Get the initial directory for the file dialog.
func (fd *fileDialog) InitialDirectory() string {
	return fd.initialDirectory
}
