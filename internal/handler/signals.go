package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sovera-core-api/internal/repository"
)

type SignalHandler struct {
	signalRepo *repository.SignalRepository
}

func NewSignalHandler(signalRepo *repository.SignalRepository) *SignalHandler {
	return &SignalHandler{signalRepo: signalRepo}
}

// ListSignals returns a paginated feed of corporate intelligence signals.
func (h *SignalHandler) ListSignals(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	minIntent, _ := strconv.Atoi(c.Query("min_intent", "70"))
	industry := c.Query("industry", "")

	signals, total, err := h.signalRepo.ListSignals(c.Context(), limit, offset, minIntent, industry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "QUERY_FAILED",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": signals,
		"pagination": fiber.Map{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// MatchPrograms runs RLS-enforced cosine vector search matching corporate signal with tenant programs.
func (h *SignalHandler) MatchPrograms(c *fiber.Ctx) error {
	signalID := c.Params("id")
	orgID, ok := c.Locals("org_id").(string)
	if !ok || orgID == "" {
		// Mock org_id for unauthenticated dev queries
		orgID = "org_77123aa-8819-4c12-99a1-00123456789a"
	}

	matches, err := h.signalRepo.MatchTenantPrograms(c.Context(), orgID, signalID, 3)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "MATCHING_FAILED",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"signal_id":            signalID,
		"top_matched_programs": matches,
	})
}
