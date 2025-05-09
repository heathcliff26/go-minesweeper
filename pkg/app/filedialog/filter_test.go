package filedialog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileFiltersExtensions(t *testing.T) {
	filters := FileFilters{
		{"Text Files", []string{".txt", ".md"}},
		{"Image Files", []string{".png", ".jpg"}},
		{"Nothing", []string{}},
		{"Nil", nil},
	}
	expected := []string{".txt", ".md", ".png", ".jpg"}

	extensions := filters.Extensions()
	assert.Equal(t, expected, extensions, "Should convert to extensions correctly")
}
