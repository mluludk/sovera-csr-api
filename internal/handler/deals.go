package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"sovera-core-api/internal/repository"
	"sovera-core-api/internal/service/ai"
	"sovera-core-api/internal/service/exporter"
)

type DealHandler struct {
	dealRepo      *repository.DealRepository
	programRepo   *repository.ProgramRepository
	signalRepo    *repository.SignalRepository
	geminiService *ai.GeminiService
	exporter      *exporter.DocumentExporter
}

func NewDealHandler(
	dealRepo *repository.DealRepository,
	programRepo *repository.ProgramRepository,
	signalRepo *repository.SignalRepository,
	geminiService *ai.GeminiService,
	exporter *exporter.DocumentExporter,
) *DealHandler {
	return &DealHandler{
		dealRepo:      dealRepo,
		programRepo:   programRepo,
		signalRepo:    signalRepo,
		geminiService: geminiService,
		exporter:      exporter,
	}
}

type CreateDealPayload struct {
	SignalID        string  `json:"signal_id"`
	CompanyName     string  `json:"company_name"`
	TargetProgramID string  `json:"target_program_id"`
	EstimatedValue  float64 `json:"estimated_value"`
}

type UpdateStagePayload struct {
	DealStage string `json:"deal_stage"`
}

type GeneratePitchPayload struct {
	Tone        string `json:"tone"`
	CustomNotes string `json:"custom_notes"`
}

func (h *DealHandler) ListDeals(c *fiber.Ctx) error {
	orgID, ok := c.Locals("org_id").(string)
	if !ok || orgID == "" {
		orgID = "org_77123aa-8819-4c12-99a1-00123456789a"
	}

	deals, err := h.dealRepo.ListDeals(c.Context(), orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "QUERY_FAILED",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": deals,
	})
}

func (h *DealHandler) CreateDeal(c *fiber.Ctx) error {
	orgID, ok := c.Locals("org_id").(string)
	if !ok || orgID == "" {
		orgID = "org_77123aa-8819-4c12-99a1-00123456789a"
	}

	var payload CreateDealPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_PAYLOAD",
			"message": "Failed to parse deal payload",
		})
	}

	if payload.CompanyName == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"error":   "MISSING_COMPANY_NAME",
			"message": "company_name is required",
		})
	}

	deal, err := h.dealRepo.CreateDeal(
		c.Context(), orgID, payload.SignalID, payload.CompanyName,
		payload.TargetProgramID, payload.EstimatedValue,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "CREATE_FAILED",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(deal)
}

func (h *DealHandler) UpdateStage(c *fiber.Ctx) error {
	dealID := c.Params("id")
	orgID, ok := c.Locals("org_id").(string)
	if !ok || orgID == "" {
		orgID = "org_77123aa-8819-4c12-99a1-00123456789a"
	}

	var payload UpdateStagePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "INVALID_PAYLOAD",
			"message": "Failed to parse stage payload",
		})
	}

	deal, err := h.dealRepo.UpdateDealStage(c.Context(), orgID, dealID, payload.DealStage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "UPDATE_FAILED",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(deal)
}

func (h *DealHandler) GeneratePitch(c *fiber.Ctx) error {
	dealID := c.Params("id")
	orgID, ok := c.Locals("org_id").(string)
	if !ok || orgID == "" {
		orgID = "org_77123aa-8819-4c12-99a1-00123456789a"
	}

	var payload GeneratePitchPayload
	_ = c.BodyParser(&payload)

	deal, err := h.dealRepo.GetDealByID(c.Context(), orgID, dealID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "DEAL_NOT_FOUND",
			"message": "Deal record not found",
		})
	}

	programTitle := "Program Beasiswa Generasi Digital 3T"
	programDesc := "Program penyediaan sarana komputer dan beasiswa digital"
	if deal.TargetProgramID != "" {
		if prog, err := h.programRepo.GetProgramByID(c.Context(), orgID, deal.TargetProgramID); err == nil {
			programTitle = prog.Title
			programDesc = prog.Description
		}
	}

	result, err := h.geminiService.GeneratePitchStrategy(
		c.Context(), dealID, deal.CompanyName, "Alokasi TJSL Pendidikan Digital",
		programTitle, programDesc, payload.Tone, payload.CustomNotes,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "AI_GENERATION_FAILED",
			"message": err.Error(),
		})
	}

	// Persist pitch strategy outputs to deal_pipelines record
	_ = h.dealRepo.UpdateDealPitch(c.Context(), orgID, dealID, result.Icebreaker, result.ProposalMarkdown)

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *DealHandler) ExportProposal(c *fiber.Ctx) error {
	dealID := c.Params("id")
	orgID, ok := c.Locals("org_id").(string)
	if !ok || orgID == "" {
		orgID = "org_77123aa-8819-4c12-99a1-00123456789a"
	}

	format := c.Query("format", "docx")

	deal, err := h.dealRepo.GetDealByID(c.Context(), orgID, dealID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "DEAL_NOT_FOUND",
			"message": "Deal record not found",
		})
	}

	content := deal.GeneratedProposal
	if content == "" {
		content = fmt.Sprintf("# PROPOSAL KEMITRAAN STRATEGIS\n\n## Korporasi: %s\n\nRingkasan draf proposal kemitraan institusional.", deal.CompanyName)
	}

	if format == "pdf" {
		pdfBytes, filename, err := h.exporter.GeneratePDF(deal.CompanyName, "Proposal Kemitraan", content)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		return c.Send(pdfBytes)
	}

	// Default format DOCX
	docxBytes, filename, err := h.exporter.GenerateDOCX(deal.CompanyName, "Proposal Kemitraan", content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return c.Send(docxBytes)
}
