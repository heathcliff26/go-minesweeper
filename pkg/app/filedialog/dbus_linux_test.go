//go:build linux

package filedialog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertFilters(t *testing.T) {
	filter := FileFilters{
		{"Text Files", []string{".txt", ".md"}},
		{"Image Files", []string{".png", ".jpg"}},
		{"Nothing", []string{}},
		{"Nil", nil},
	}
	expected := []freedesktopFilter{
		{
			Name: "Text Files",
			Rules: []freedesktopFilterRule{
				{Pattern: "*.txt"},
				{Pattern: "*.md"},
			},
		},
		{
			Name: "Image Files",
			Rules: []freedesktopFilterRule{
				{Pattern: "*.png"},
				{Pattern: "*.jpg"},
			},
		},
		{
			Name: "Nothing",
		},
		{
			Name: "Nil",
		},
	}

	converted := convertFilters(filter)

	assert.Equal(t, expected, converted, "Should convert filters correctly")
}
