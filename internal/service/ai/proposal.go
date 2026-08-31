package ai

import (
	"context"
	"fmt"
)

type PitchStrategyResult struct {
	DealID           string   `json:"deal_id"`
	Icebreaker       string   `json:"icebreaker"`
	PitchDeckOutline []string `json:"pitch_deck_outline"`
	ProposalMarkdown string   `json:"proposal_markdown"`
}

// GeneratePitchStrategy orchestrates Gemini LLM to craft executive outreach outputs.
func (s *GeminiService) GeneratePitchStrategy(ctx context.Context, dealID, companyName, CSRSummary, programTitle, programDesc, tone, customNotes string) (*PitchStrategyResult, error) {
	// Fallback/Mock synthesis logic for instant strategy generation
	icebreaker := fmt.Sprintf(
		"Yth. Pimpinan TJSL & CSR %s,\n\nMenyikapi inisiatif luar biasa korporasi dalam fokus keberlanjutan (%s), kami dari lembaga bermaksud mengajukan kolaborasi strategis melalui %s.\n\n%s\n\nBesar harapan kami untuk dapat menjadwalkan sesi audiensi singkat.",
		companyName, CSRSummary, programTitle, customNotes,
	)

	pitchOutline := []string{
		"Slide 1: Urgensi Kesenjangan & Tantangan Akses Sektor Target",
		fmt.Sprintf("Slide 2: Solusi Kemitraan Strategis via %s", programTitle),
		"Slide 3: Rencana Target Dampak (SDGs & Indikator Penerima Manfaat)",
		"Slide 4: Skema Pengelolaan Anggaran & Transparansi Laporan",
		"Slide 5: Profil Lembaga & Rekam Jejak Kemitraan Korporasi",
	}

	proposalMarkdown := fmt.Sprintf(
		"# PROPOSAL KEMITRAAN STRATEGIS\n\n## Target Korporasi: %s\n## Program Lembaga: %s\n\n### 1. Ringkasan Eksekutif\n%s\n\n### 2. Deskripsi Program\n%s\n\n### 3. Catatan Khusus & Penyesuaian Strategis\n%s",
		companyName, programTitle, CSRSummary, programDesc, customNotes,
	)

	return &PitchStrategyResult{
		DealID:           dealID,
		Icebreaker:       icebreaker,
		PitchDeckOutline: pitchOutline,
		ProposalMarkdown: proposalMarkdown,
	}, nil
}
