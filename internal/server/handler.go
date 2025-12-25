package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"dynamic-pricing-tool-ru/internal/logger"
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
		logger.L.Warn("Invalid request format",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.Mapping == nil || req.Data == nil {
		logger.L.Warn("Missing required fields")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required fields",
		})
		return
	}

	results, err := h.processor.ProcessRequest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	analysis := processor.AnalyzeResults(results)

	response := map[string]interface{}{
		"status":   "success",
		"analysis": analysis,
		"results":  results,
		"count":    len(results),
	}

	// Если нужны только упрощенные данные
	includeRaw := c.Query("includeRaw") == "true"
	if !includeRaw {
		// Удаляем сырые данные из ответа
		for i := range results {
			results[i].GetchipsRaw = nil
			results[i].EfindRaw = nil
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "part-api-processor",
		"apis":    []string{"Getchips", "Efind"},
	})
}
