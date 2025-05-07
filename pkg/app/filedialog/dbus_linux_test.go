//go:build linux

package filedialog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertFilters(t *testing.T) {
	filter := FileFilter{
		"Text Files":  []string{".txt", ".md"},
		"Image Files": []string{".png", ".jpg"},
		"Nothing":     []string{},
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
	}

	converted := convertFilters(filter)

	assert.ElementsMatch(t, expected, converted, "Should convert filters correctly")
}
