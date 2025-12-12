package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dynamic-pricing-tool-ru/internal/processor"
	"dynamic-pricing-tool-ru/internal/types"
)

type Handler struct {
	processor *processor.Processor
}

func NewHandler(proc *processor.Processor) *Handler {
	return &Handler{
		processor: proc,
	}
}

func (h *Handler) HandleProcess(c *gin.Context) {
	var req types.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate request
	if req.Mapping == nil || req.Data == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Process request
	results, err := h.processor.ProcessRequest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format response
	response := map[string]interface{}{
		"status":  "success",
		"count":   len(results),
		"results": results,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "part-api-processor",
	})
}
