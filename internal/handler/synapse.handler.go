package handler

import (
	"errors"
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type SynapseHandler struct {
	synapseService service.SynapseService
}

func NewSynapseHandler(synapseService service.SynapseService) *SynapseHandler {
	return &SynapseHandler{synapseService: synapseService}
}

func (h *SynapseHandler) FindSynapse(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	synapseULID := req.ULID

	if synapseULID == "" {
		c.JSON(http.StatusBadRequest, "invalid synapse ID")
		return errors.New("invalid synapse ID")
	}

	synapse, err := h.synapseService.GetSynapseByULID(c.Request().Context(), synapseULID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, synapse)
	return nil
}

func (h *SynapseHandler) SearchSynapses(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	synapses, err := h.synapseService.SearchSynapses(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, synapses)
	return nil
}

func (h *SynapseHandler) CountSynapses(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	req.Count = true

	count, err := h.synapseService.CountSynapses(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, count)
	return nil
}
