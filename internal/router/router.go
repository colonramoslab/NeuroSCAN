package router

import (
	"neuroscan/internal/handler"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, neuronHandler *handler.NeuronHandler, contactHandler *handler.ContactHandler, synapseHandler *handler.SynapseHandler, cphateHandler *handler.CphateHandler, nerveringHandler *handler.NerveRingHandler, scaleHandler *handler.ScaleHandler, promoterHandler *handler.PromoterHandler, developmentalStageHandler *handler.DevelopmentalStageHandler, videoHandler *handler.VideoHandler) *echo.Echo {
	e.GET("/neurons", neuronHandler.SearchNeurons)
	e.GET("/neurons/:ulid", neuronHandler.FindNeuronByULID)
	e.GET("/neurons/:timepoint/:uid", neuronHandler.FindNeuronByUID)
	e.GET("/neurons/count", neuronHandler.CountNeurons)

	e.GET("/contacts", contactHandler.SearchContacts)
	e.GET("/contacts/:ulid", contactHandler.FindContactByULID)
	e.GET("/contacts/:timepoint/:uid", contactHandler.FindContactByUID)
	e.GET("/contacts/count", contactHandler.CountContacts)

	e.GET("/synapses", synapseHandler.SearchSynapses)
	e.GET("/synapses/:ulid", synapseHandler.FindSynapseByULID)
	e.GET("/synapses/:timepoint/:uid", synapseHandler.FindSynapseByUID)
	e.GET("/synapses/count", synapseHandler.CountSynapses)

	e.GET("/cphates", cphateHandler.CphateByTimepoint)
	e.GET("/cphates/count", cphateHandler.CountCphates)

	e.GET("/nerve-rings", nerveringHandler.NerveRingByTimepoint)

	e.GET("/scales", scaleHandler.ScaleByTimepoint)

	e.GET("/promoters", promoterHandler.SearchPromoters)

	e.GET("/developmental-stages", developmentalStageHandler.SearchDevelopmentalStages)
	e.GET("/developmental-stages/count", developmentalStageHandler.CountDevelopmentalStages)

	e.POST("/videos/webmtomp4", videoHandler.UploadWebm)
	e.GET("/videos/status/:uuid", videoHandler.UploadStatus)
	e.GET("/videos/download/:filename", videoHandler.DownloadMP4)

	return e
}
