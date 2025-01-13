package handler

import (
	"net/http"

	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type NeuronHandler struct {
	neuronService service.NeuronService
}

func NewNeuronHandler(neuronService service.NeuronService) *NeuronHandler {
	return &NeuronHandler{neuronService: neuronService}
}

// GetNeurons returns all neurons
func (h *NeuronHandler) GetAllNeurons(c echo.Context) error {
	var req APIRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	neurons, err := h.neuronService.GetAllNeurons(c.Request().Context(), req.Timepoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, neurons)
	return nil
}

