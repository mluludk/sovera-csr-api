package normalizer

import (
	"testing"
)

func TestCleanText(t *testing.T) {
	norm := NewNormalizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "html tags removal",
			input:    "<h1>Title</h1><p>Hello <b>world</b>!</p>",
			expected: "Title Hello world !",
		},
		{
			name:     "excessive whitespace",
			input:    "   Multiple    spaces   and\n\n\n\nmultiple   newlines.   ",
			expected: "Multiple spaces and\n\nmultiple newlines.",
		},
		{
			name:     "carriage returns",
			input:    "Line 1\r\nLine 2\rLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := norm.CleanText(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSelectBestContent(t *testing.T) {
	norm := NewNormalizer()

	t.Run("uses raw text when present", func(t *testing.T) {
		got := norm.SelectBestContent("<p>Raw Content</p>", "Markdown Content")
		if got != "Raw Content" {
			t.Errorf("expected 'Raw Content', got %q", got)
		}
	})

	t.Run("falls back to markdown content when raw text empty", func(t *testing.T) {
		got := norm.SelectBestContent("", "# Markdown Content")
		if got != "# Markdown Content" {
			t.Errorf("expected '# Markdown Content', got %q", got)
		}
	})
}

func TestGenerateContentHash(t *testing.T) {
	norm := NewNormalizer()

	hash1 := norm.GenerateContentHash("Sample Text 123")
	hash2 := norm.GenerateContentHash("Sample Text 123")
	hash3 := norm.GenerateContentHash("Different Text")

	if hash1 != hash2 {
		t.Errorf("expected identical hashes for same content, got %s and %s", hash1, hash2)
	}
	if hash1 == hash3 {
		t.Errorf("expected different hashes for different content, got same: %s", hash1)
	}
	if len(hash1) != 64 {
		t.Errorf("expected 64 character SHA-256 hex string, got length %d", len(hash1))
	}
}

func TestPrepareEmbeddingText(t *testing.T) {
	norm := NewNormalizer()

	got := norm.PrepareEmbeddingText("PT Telkom Tbk", "Telecommunication", "Digital Education", "Beasiswa 3T")
	expected := "PT Telkom Tbk Telecommunication Digital Education Beasiswa 3T"

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
