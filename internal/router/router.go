package router

import (
	"neuroscan/internal/handler"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, neuronHandler *handler.NeuronHandler, contactHandler *handler.ContactHandler) *echo.Echo {
	e.GET("/neurons", neuronHandler.SearchNeurons)
	e.GET("/neurons/count", neuronHandler.CountNeurons)

	e.GET("/contacts", contactHandler.SearchContacts)
	e.GET("/contacts/count", contactHandler.CountContacts)

	return e
}
