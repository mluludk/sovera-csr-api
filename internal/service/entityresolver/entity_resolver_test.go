package entityresolver

import (
	"context"
	"testing"
)

func TestCleanCompanyName(t *testing.T) {
	resolver := NewEntityResolver(nil)

	tests := []struct {
		input    string
		expected string
	}{
		{"PT Telkom Indonesia (Persero) Tbk.", "Telkom Indonesia"},
		{"PT. Maju Bersama Tbk", "Maju Bersama"},
		{"Bank Central Asia (BCA) Tbk", "Bank Central Asia BCA"},
		{"Astra International Group Inc.", "Astra International"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := resolver.CleanCompanyName(tt.input)
			if got != tt.expected {
				t.Errorf("CleanCompanyName(%q) = %q; expected %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	resolver := NewEntityResolver(nil)

	tests := []struct {
		input    string
		expected string
	}{
		{"PT Telkom Indonesia (Persero) Tbk.", "telkom-indonesia"},
		{"PT. Maju Bersama Tbk", "maju-bersama"},
		{"Astra International", "astra-international"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := resolver.GenerateSlug(tt.input)
			if got != tt.expected {
				t.Errorf("GenerateSlug(%q) = %q; expected %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestResolveCompany_Mock(t *testing.T) {
	resolver := NewEntityResolver(nil)

	res, err := resolver.ResolveCompany(context.Background(), "PT Telkom Indonesia Tbk", "Telecommunication")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Slug != "telkom-indonesia" {
		t.Errorf("expected slug 'telkom-indonesia', got %q", res.Slug)
	}
	if res.CompanyID == nil {
		t.Errorf("expected non-nil CompanyID")
	}
}
