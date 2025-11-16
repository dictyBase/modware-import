package stockcenter

import (
	"testing"

	E "github.com/IBM/fp-go/either"
	O "github.com/IBM/fp-go/option"
	"github.com/stretchr/testify/require"
)

func TestParseCommaSeparatedField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected O.Option[[]string]
	}{
		{
			name:     "empty string returns None",
			input:    "",
			expected: O.None[[]string](),
		},
		{
			name:     "whitespace only returns None",
			input:    "   ",
			expected: O.None[[]string](),
		},
		{
			name:     "single gene",
			input:    "abpC",
			expected: O.Some([]string{"abpC"}),
		},
		{
			name:     "multiple genes with spaces",
			input:    "abpC, gtaC, H2Bv3",
			expected: O.Some([]string{"abpC", "gtaC", "H2Bv3"}),
		},
		{
			name:     "genes with extra whitespace",
			input:    "  abpC ,  gtaC  ",
			expected: O.Some([]string{"abpC", "gtaC"}),
		},
		{
			name:     "single gene with trailing comma",
			input:    "abpC,",
			expected: O.Some([]string{"abpC"}),
		},
		{
			name:     "empty elements filtered out",
			input:    "abpC,,gtaC",
			expected: O.Some([]string{"abpC", "gtaC"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommaSeparatedField(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestParsePublicationField(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected O.Option[[]string]
	}{
		{
			name:     "empty string returns None",
			input:    "",
			expected: O.None[[]string](),
		},
		{
			name:     "whitespace only returns None",
			input:    "   ",
			expected: O.None[[]string](),
		},
		{
			name:     "single PMID",
			input:    "32232356",
			expected: O.Some([]string{"32232356"}),
		},
		{
			name:     "PMID with whitespace",
			input:    "  32232356  ",
			expected: O.Some([]string{"32232356"}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePublicationField(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestHasValidRecordLength(t *testing.T) {
	tests := []struct {
		name     string
		record   []string
		expected bool
	}{
		{
			name:     "valid 7 fields",
			record:   []string{"a", "b", "c", "d", "e", "f", "g"},
			expected: true,
		},
		{
			name:     "invalid - too few fields",
			record:   []string{"a", "b"},
			expected: false,
		},
		{
			name:     "invalid - too many fields",
			record:   []string{"a", "b", "c", "d", "e", "f", "g", "h"},
			expected: false,
		},
		{
			name:     "invalid - empty",
			record:   []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasValidRecordLength(tt.record)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildPlasmid(t *testing.T) {
	tests := []struct {
		name      string
		record    []string
		userEmail string
		cvterm    string
		validate  func(*testing.T, *GoldenBraidPlasmid)
	}{
		{
			name: "valid record with all fields",
			record: []string{
				"pDGB_A1",          // Name
				"pDGB_A1",          // Synonym (ignored)
				"Peter Kundert",    // Depositor (ignored)
				"abpC, gtaC",       // Genes
				"GB backbone",      // Keywords (ignored)
				"Test description", // Summary
				"32232356",         // PMID
			},
			userEmail: "test@example.com",
			cvterm:    "GB vector",
			validate: func(t *testing.T, p *GoldenBraidPlasmid) {
				require.Equal(t, "pDGB_A1", p.Name)
				require.Equal(t, "Test description", p.Summary)
				require.Equal(t, "test@example.com", p.User)
				require.Equal(t, "GB vector", p.PlasmidType)

				// Check genes
				require.True(t, O.IsSome(p.Genes))
				genes := O.GetOrElse(func() []string { return []string{} })(p.Genes)
				require.Equal(t, []string{"abpC", "gtaC"}, genes)

				// Check publications
				require.True(t, O.IsSome(p.Publications))
				pubs := O.GetOrElse(func() []string { return []string{} })(p.Publications)
				require.Equal(t, []string{"32232356"}, pubs)
			},
		},
		{
			name: "valid record with empty genes",
			record: []string{
				"pDGB_A1",
				"pDGB_A1",
				"Peter Kundert",
				"", // Empty genes
				"GB backbone",
				"Test description",
				"32232356",
			},
			userEmail: "test@example.com",
			cvterm:    "GB vector",
			validate: func(t *testing.T, p *GoldenBraidPlasmid) {
				require.True(t, O.IsNone(p.Genes))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plasmid := BuildPlasmid(tt.userEmail, tt.cvterm)(tt.record)
			require.NotNil(t, plasmid)
			if tt.validate != nil {
				tt.validate(t, plasmid)
			}
		})
	}
}

func TestValidatePlasmid(t *testing.T) {
	tests := []struct {
		name      string
		plasmid   *GoldenBraidPlasmid
		shouldErr bool
		errorMsg  string
	}{
		{
			name: "valid plasmid",
			plasmid: &GoldenBraidPlasmid{
				Name:    "pDGB_A1",
				Summary: "Test plasmid",
				User:    "test@example.com",
			},
			shouldErr: false,
		},
		{
			name: "invalid name - no 'p' prefix",
			plasmid: &GoldenBraidPlasmid{
				Name:    "DGB_A1",
				Summary: "Test",
				User:    "test@example.com",
			},
			shouldErr: true,
			errorMsg:  "must start with 'p'",
		},
		{
			name: "invalid name - empty",
			plasmid: &GoldenBraidPlasmid{
				Name:    "",
				Summary: "Test",
				User:    "test@example.com",
			},
			shouldErr: true,
			errorMsg:  "must start with 'p'",
		},
		{
			name: "invalid summary - empty",
			plasmid: &GoldenBraidPlasmid{
				Name:    "pDGB_A1",
				Summary: "",
				User:    "test@example.com",
			},
			shouldErr: true,
			errorMsg:  "empty summary",
		},
		{
			name: "invalid user - empty",
			plasmid: &GoldenBraidPlasmid{
				Name:    "pDGB_A1",
				Summary: "Test",
				User:    "",
			},
			shouldErr: true,
			errorMsg:  "no user email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePlasmid(tt.plasmid)

			if tt.shouldErr {
				require.True(t, E.IsLeft(result))
				err := E.Fold(
					func(e error) error { return e },
					func(*GoldenBraidPlasmid) error { return nil },
				)(result)
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.True(t, E.IsRight(result))
			}
		})
	}
}

//nolint:funlen // Table-driven test with multiple predicates
func TestValidationPredicates(t *testing.T) {
	t.Run("hasValidName", func(t *testing.T) {
		tests := []struct {
			name     string
			plasmid  *GoldenBraidPlasmid
			expected bool
		}{
			{
				name:     "valid name with p prefix",
				plasmid:  &GoldenBraidPlasmid{Name: "pDGB_A1"},
				expected: true,
			},
			{
				name:     "invalid name without p prefix",
				plasmid:  &GoldenBraidPlasmid{Name: "DGB_A1"},
				expected: false,
			},
			{
				name:     "empty name",
				plasmid:  &GoldenBraidPlasmid{Name: ""},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := hasValidName(tt.plasmid)
				require.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("hasValidSummary", func(t *testing.T) {
		tests := []struct {
			name     string
			plasmid  *GoldenBraidPlasmid
			expected bool
		}{
			{
				name:     "valid summary",
				plasmid:  &GoldenBraidPlasmid{Summary: "Test description"},
				expected: true,
			},
			{
				name:     "empty summary",
				plasmid:  &GoldenBraidPlasmid{Summary: ""},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := hasValidSummary(tt.plasmid)
				require.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("hasValidUser", func(t *testing.T) {
		tests := []struct {
			name     string
			plasmid  *GoldenBraidPlasmid
			expected bool
		}{
			{
				name:     "valid user",
				plasmid:  &GoldenBraidPlasmid{User: "test@example.com"},
				expected: true,
			},
			{
				name:     "empty user",
				plasmid:  &GoldenBraidPlasmid{User: ""},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := hasValidUser(tt.plasmid)
				require.Equal(t, tt.expected, result)
			})
		}
	})
}
