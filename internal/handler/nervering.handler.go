package handler

import (
	"errors"
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type NerveRingHandler struct {
	nerveringService service.NerveRingService
}

func NewNerveRingHandler(nerveringService service.NerveRingService) *NerveRingHandler {
	return &NerveRingHandler{nerveringService: nerveringService}
}

func (h *NerveRingHandler) NerveRingByTimepoint(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	if req.Timepoint == nil {
		c.JSON(http.StatusBadRequest, errors.New("timepoint is required"))
		return errors.New("timepoint is required")
	}

	nervering, err := h.nerveringService.GetNerveRingByTimepoint(c.Request().Context(), *req.Timepoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, nervering)
	return nil
}

