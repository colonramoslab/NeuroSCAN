package handler

import (
	"net/http"

	"neuroscan/internal/domain"
	"neuroscan/internal/service"

	"github.com/labstack/echo/v4"
)

type ContactHandler struct {
	contactService service.ContactService
}

func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

func (h *ContactHandler) SearchContacts(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	contacts, err := h.contactService.SearchContacts(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}

	c.JSON(http.StatusOK, contacts)
	return nil
}

func (h *ContactHandler) CountContacts(c echo.Context) error {
	var req domain.APIV1Request

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	req.Count = true

	count, err := h.contactService.CountContacts(c.Request().Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}
	c.JSON(http.StatusOK, count)
	return nil
}
