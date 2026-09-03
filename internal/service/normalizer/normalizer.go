package normalizer

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

var (
	htmlTagRegex       = regexp.MustCompile(`(?i)<[^>]*>`)
	multiSpaceRegex    = regexp.MustCompile(`[ \t]+`)
	multiNewlineRegex  = regexp.MustCompile(`\n{3,}`)
	controlCharRegex   = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F]`)
)

type Normalizer struct{}

func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

// CleanText removes HTML tags, non-printable control characters, and normalizes excessive whitespace.
func (n *Normalizer) CleanText(input string) string {
	if input == "" {
		return ""
	}

	// 1. Remove control characters
	cleaned := controlCharRegex.ReplaceAllString(input, "")

	// 2. Remove HTML tags
	cleaned = htmlTagRegex.ReplaceAllString(cleaned, " ")

	// 3. Replace Windows carriage returns
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "\n")

	// 4. Compress horizontal spaces
	cleaned = multiSpaceRegex.ReplaceAllString(cleaned, " ")

	// 5. Compress multiple blank lines (max 2 newlines)
	cleaned = multiNewlineRegex.ReplaceAllString(cleaned, "\n\n")

	return strings.TrimSpace(cleaned)
}

// SelectBestContent returns the cleaned primary text, falling back from rawText to markdownContent.
func (n *Normalizer) SelectBestContent(rawText, markdownContent string) string {
	cleanedRaw := n.CleanText(rawText)
	if cleanedRaw != "" {
		return cleanedRaw
	}
	return n.CleanText(markdownContent)
}

// GenerateContentHash creates a SHA-256 hex hash for deduplication.
func (n *Normalizer) GenerateContentHash(text string) string {
	cleaned := n.CleanText(text)
	hash := sha256.Sum256([]byte(cleaned))
	return hex.EncodeToString(hash[:])
}

// PrepareEmbeddingText formats and normalizes composite metadata into a string suitable for 1536-dim vector embedding.
func (n *Normalizer) PrepareEmbeddingText(companyName, industry, pillar, summary string) string {
	parts := []string{
		n.CleanText(companyName),
		n.CleanText(industry),
		n.CleanText(pillar),
		n.CleanText(summary),
	}

	nonEmpty := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}

	return strings.Join(nonEmpty, " ")
}
