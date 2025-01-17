package handler

import (
	"errors"
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type CphateHandler struct {
	cphateService service.CphateService
}

func NewCphateHandler(cphateService service.CphateService) *CphateHandler {
	return &CphateHandler{cphateService: cphateService}
}

func (h *CphateHandler) CphateByTimepoint(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	if req.Timepoint == nil {
		c.JSON(http.StatusBadRequest, errors.New("timepoint is required"))
		return errors.New("timepoint is required")
	}

	cphate, err := h.cphateService.GetCphateByTimepoint(c.Request().Context(), *req.Timepoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, cphate)
	return nil
}
