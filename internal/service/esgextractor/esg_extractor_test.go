package esgextractor

import (
	"context"
	"testing"

	"sovera-core-api/internal/service/ai"
	"sovera-core-api/internal/service/entityresolver"
	"sovera-core-api/internal/service/normalizer"
)

func TestProcessESGExtraction_Mock(t *testing.T) {
	gemini := ai.NewGeminiService("") // Uses mock
	norm := normalizer.NewNormalizer()
	resolver := entityresolver.NewEntityResolver(nil)

	extractor := NewESGExtractor(gemini, nil, norm, resolver)

	compID := "comp_12345678"
	rawText := "PT Telkom Indonesia Laporan Keberlanjutan 2024. Skor ESG 84.5 dengan rating AA."

	profile, err := extractor.ProcessESGExtraction(
		context.Background(),
		rawText, "",
		"PT Telkom Indonesia Tbk", "Telecommunication",
		&compID,
	)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if profile.CompanyID != compID {
		t.Errorf("expected CompanyID %q, got %q", compID, profile.CompanyID)
	}

	if profile.ESGRating == nil || *profile.ESGRating != "AA" {
		t.Errorf("expected ESG rating 'AA', got %v", profile.ESGRating)
	}
}
