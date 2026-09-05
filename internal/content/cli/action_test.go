package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNameAndNamespace(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedName      string
		expectedNamespace string
	}{
		{
			name:              "frontpage content",
			input:             "some/path/dfp-about.json",
			expectedName:      "about",
			expectedNamespace: "frontpage",
		},
		{
			name:              "stockcenter content",
			input:             "some/path/dsc-information.json",
			expectedName:      "information",
			expectedNamespace: "stockcenter",
		},
		{
			name:              "news content",
			input:             "some/path/news-256cd371-8710-462a-9f7c-d34774526c8f.json",
			expectedName:      "256cd371-8710-462a-9f7c-d34774526c8f",
			expectedNamespace: "news",
		},
		{
			name:              "no directory prefix (flat key)",
			input:             "dfp-contact.json",
			expectedName:      "contact",
			expectedNamespace: "frontpage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, namespace := nameAndNamespace(tt.input)
			require.Equal(t, tt.expectedName, name)
			require.Equal(t, tt.expectedNamespace, namespace)
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "frontpage art",
			expected: "frontpage-art",
		},
		{
			input:    "already-slugified",
			expected: "already-slugified",
		},
		{
			input:    "  leading and trailing spaces  ",
			expected: "leading-and-trailing-spaces",
		},
		{
			input:    "special!@#characters",
			expected: "special-characters",
		},
		{
			input:    "multiple   spaces",
			expected: "multiple-spaces",
		},
		{
			input:    "news 256cd371-8710-462a-9f7c-d34774526c8f",
			expected: "news-256cd371-8710-462a-9f7c-d34774526c8f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			require.Equal(t, tt.expected, Slugify(tt.input))
		})
	}
}
