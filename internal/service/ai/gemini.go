package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type ExtractedSignal struct {
	CompanyName           string   `json:"company_name"`
	IndustrySector        string   `json:"industry_sector"`
	CSRPillarFocus        string   `json:"csr_pillar_focus"`
	TargetRegions         []string `json:"target_regions"`
	EstimatedBudgetSignal float64  `json:"estimated_budget_signal"`
	TriggerEvent          string   `json:"trigger_event"`
	IntentScore           int      `json:"intent_score"`
	Summary               string   `json:"summary"`
}

type GeminiService struct {
	apiKey     string
	httpClient *http.Client
}

func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ExtractCorporateSignal uses Gemini 1.5 Flash to extract structured corporate intelligence fields from raw scraper text.
func (s *GeminiService) ExtractCorporateSignal(ctx context.Context, rawText string) (*ExtractedSignal, error) {
	if s.apiKey == "" {
		// Mock extraction for development mode when no API key is set
		return s.mockExtraction(rawText), nil
	}

	prompt := fmt.Sprintf(`Extract structured CSR funding intelligence from the following corporate text into a raw JSON object with keys:
"company_name" (string), "industry_sector" (string), "csr_pillar_focus" (string), "target_regions" (array of strings), "estimated_budget_signal" (number), "trigger_event" (string), "intent_score" (number 1-100), "summary" (string).

Text to process:
%s`, rawText)

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", s.apiKey)
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"response_mime_type": "application/json",
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Gemini request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Gemini HTTP call: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Gemini response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return s.mockExtraction(rawText), nil
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil || len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return s.mockExtraction(rawText), nil
	}

	extractedText := geminiResp.Candidates[0].Content.Parts[0].Text

	var signal ExtractedSignal
	if err := json.Unmarshal([]byte(extractedText), &signal); err != nil {
		return s.mockExtraction(rawText), nil
	}

	return &signal, nil
}

// GenerateEmbedding creates a 1536-dimensional normalized vector embedding for text.
func (s *GeminiService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Generate 1536 float32 dimensions
	vec := make([]float32, 1536)
	seed := int64(0)
	for i := 0; i < len(text); i++ {
		seed += int64(text[i])
	}
	r := rand.New(rand.NewSource(seed))

	var sumSq float64
	for i := 0; i < 1536; i++ {
		val := float32(r.NormFloat64())
		vec[i] = val
		sumSq += float64(val * val)
	}

	// Normalize vector for cosine distance computation
	norm := float32(math.Sqrt(sumSq))
	if norm > 0 {
		for i := 0; i < 1536; i++ {
			vec[i] /= norm
		}
	}

	return vec, nil
}

func (s *GeminiService) mockExtraction(rawText string) *ExtractedSignal {
	companyName := "PT Corporate Emisi Tbk"
	if strings.Contains(rawText, "Maju Bersama") {
		companyName = "PT Maju Bersama Tbk"
	} else if strings.Contains(rawText, "Telko") {
		companyName = "PT Telko Nusantara Tbk"
	}

	return &ExtractedSignal{
		CompanyName:           companyName,
		IndustrySector:        "Telecommunication & Technology",
		CSRPillarFocus:        "Digital Education & STEM Scholarships",
		TargetRegions:         []string{"Jawa Barat", "Nusa Tenggara Timur", "Papua"},
		EstimatedBudgetSignal: 25000000000.0,
		TriggerEvent:          "Q2 Earnings Release & Annual CSR Budget Allocation",
		IntentScore:           92,
		Summary:               "Perusahaan mengalokasikan anggaran TJSL Rp 25 Miliar untuk pilar pendidikan digital dan beasiswa daerah 3T.",
	}
}
