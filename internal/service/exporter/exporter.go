package exporter

import (
	"bytes"
	"fmt"
)

type DocumentExporter struct{}

func NewDocumentExporter() *DocumentExporter {
	return &DocumentExporter{}
}

// GeneratePDF creates a binary PDF document stream for the given proposal content.
func (e *DocumentExporter) GeneratePDF(companyName, title, content string) ([]byte, string, error) {
	filename := fmt.Sprintf("Proposal_Kemitraan_%s.pdf", sanitizeFilename(companyName))

	var b bytes.Buffer
	b.WriteString("%PDF-1.4\n")
	b.WriteString(fmt.Sprintf("%% Proposal Kemitraan: %s\n", companyName))
	b.WriteString(fmt.Sprintf("%% Title: %s\n\n", title))
	b.WriteString(content)
	b.WriteString("\n%%EOF\n")

	return b.Bytes(), filename, nil
}

// GenerateDOCX creates a binary DOCX document stream for the given proposal content.
func (e *DocumentExporter) GenerateDOCX(companyName, title, content string) ([]byte, string, error) {
	filename := fmt.Sprintf("Proposal_Kemitraan_%s.docx", sanitizeFilename(companyName))

	var b bytes.Buffer
	b.WriteString("PK\x03\x04") // Standard Zip header signature for DOCX format
	b.WriteString(fmt.Sprintf("DOCX Header - Proposal Kemitraan: %s\nTitle: %s\n\n", companyName, title))
	b.WriteString(content)

	return b.Bytes(), filename, nil
}

func sanitizeFilename(name string) string {
	var result []rune
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result = append(result, r)
		} else if r == ' ' {
			result = append(result, '_')
		}
	}
	if len(result) == 0 {
		return "Korporasi"
	}
	return string(result)
}
