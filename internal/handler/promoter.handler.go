package handler

import (
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type PromoterHandler struct {
	promoterService service.PromoterService
}

func NewPromoterHandler(promoterService service.PromoterService) *PromoterHandler {
	return &PromoterHandler{promoterService: promoterService}
}

func (h *PromoterHandler) SearchPromoters(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	promoters, err := h.promoterService.SearchPromoters(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, promoters)
	return nil
}
