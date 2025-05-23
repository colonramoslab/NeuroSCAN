package handler

import (
	"errors"
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type NeuronHandler struct {
	neuronService service.NeuronService
}

func NewNeuronHandler(neuronService service.NeuronService) *NeuronHandler {
	return &NeuronHandler{neuronService: neuronService}
}

func (h *NeuronHandler) FindNeuron(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	neuronULID := req.ULID

	if neuronULID == "" {
		c.JSON(http.StatusBadRequest, "invalid neuron ID")
		return errors.New("invalid neuron ID")
	}

	neuron, err := h.neuronService.GetNeuronByULID(c.Request().Context(), neuronULID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, neuron)
	return nil
}

func (h *NeuronHandler) SearchNeurons(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	neurons, err := h.neuronService.SearchNeurons(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, neurons)
	return nil
}

func (h *NeuronHandler) CountNeurons(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	req.Count = true

	count, err := h.neuronService.CountNeurons(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}
	c.JSON(http.StatusOK, count)
	return nil
}
