package handler

import (
	"errors"
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type ScaleHandler struct {
	scaleService service.ScaleService
}

func NewScaleHandler(scaleService service.ScaleService) *ScaleHandler {
	return &ScaleHandler{scaleService: scaleService}
}

func (h *ScaleHandler) ScaleByTimepoint(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	if req.Timepoint == nil {
		c.JSON(http.StatusBadRequest, errors.New("timepoint is required"))
		return errors.New("timepoint is required")
	}

	scale, err := h.scaleService.GetScaleByTimepoint(c.Request().Context(), *req.Timepoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, scale)
	return nil
}
