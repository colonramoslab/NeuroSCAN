package handler

import (
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type DevelopmentalStageHandler struct {
	developmentalStageService service.DevelopmentalStageService
}

func NewDevelopmentalStageHandler(developmentalStageService service.DevelopmentalStageService) *DevelopmentalStageHandler {
	return &DevelopmentalStageHandler{developmentalStageService: developmentalStageService}
}

func (h *DevelopmentalStageHandler) SearchDevelopmentalStages(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	developmentalStages, err := h.developmentalStageService.SearchDevelopmentalStages(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, developmentalStages)
	return nil
}

func (h *DevelopmentalStageHandler) CountDevelopmentalStages(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	req.Count = true

	count, err := h.developmentalStageService.CountDevelopmentalStages(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, count)
	return nil
}
