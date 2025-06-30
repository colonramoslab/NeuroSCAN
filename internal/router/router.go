package router

import (
	"neuroscan/internal/handler"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, neuronHandler *handler.NeuronHandler, contactHandler *handler.ContactHandler, synapseHandler *handler.SynapseHandler, cphateHandler *handler.CphateHandler, nerveringHandler *handler.NerveRingHandler, scaleHandler *handler.ScaleHandler, promoterHandler *handler.PromoterHandler, developmentalStageHandler *handler.DevelopmentalStageHandler) *echo.Echo {
	e.GET("/neurons", neuronHandler.SearchNeurons)
	e.GET("/neurons/:ulid", neuronHandler.FindNeuron)
	e.GET("/neurons/count", neuronHandler.CountNeurons)

	e.GET("/contacts", contactHandler.SearchContacts)
	e.GET("/contacts/:ulid", contactHandler.FindContact)
	e.GET("/contacts/count", contactHandler.CountContacts)

	e.GET("/synapses", synapseHandler.SearchSynapses)
	e.GET("/synapses/:ulid", synapseHandler.FindSynapse)
	e.GET("/synapses/count", synapseHandler.CountSynapses)

	e.GET("/cphates", cphateHandler.CphateByTimepoint)
	e.GET("/cphates/count", cphateHandler.CountCphates)

	e.GET("/nerve-rings", nerveringHandler.NerveRingByTimepoint)

	e.GET("/scales", scaleHandler.ScaleByTimepoint)

	e.GET("/promoters", promoterHandler.SearchPromoters)

	e.GET("/developmental-stages", developmentalStageHandler.SearchDevelopmentalStages)
	e.GET("/developmental-stages/count", developmentalStageHandler.CountDevelopmentalStages)

	e.POST("/webmtomp4", handler.UploadWebm)

	return e
}
