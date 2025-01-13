package router

import (
	"neuroscan/internal/handler"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, neuronHandler *handler.NeuronHandler) *echo.Echo {
	e.GET("/neurons", neuronHandler.SearchNeurons)
	e.GET("/neurons/count", neuronHandler.CountNeurons)
	
	return e
}
