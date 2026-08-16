package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"aftercare/internal/aftercare"
)

const maxRequestBytes = 1 << 20

type Handler struct {
	service *aftercare.Service
}

func NewHandler(service *aftercare.Service) http.Handler {
	if service == nil {
		panic("aftercare service is required")
	}

	handler := &Handler{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handler.health)
	mux.HandleFunc("POST /api/tickets", handler.createTicket)
	mux.HandleFunc("GET /api/tickets/{id}", handler.findTicket)
	return mux
}

func (h *Handler) health(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) createTicket(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Description *string `json:"description"`
	}

	decoder := json.NewDecoder(http.MaxBytesReader(writer, request.Body, maxRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, aftercare.ErrorView("invalid_request", "request body must be a valid ticket object"))
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeJSON(writer, http.StatusBadRequest, aftercare.ErrorView("invalid_request", "request body must contain one ticket object"))
		return
	}
	if input.Description == nil {
		writeJSON(writer, http.StatusUnprocessableEntity, aftercare.ErrorView("description_required", "description is required"))
		return
	}

	view := h.service.Submit(request.Context(), *input.Description, nil)
	if view.State == aftercare.StateError {
		writeJSON(writer, http.StatusUnprocessableEntity, view)
		return
	}
	writeJSON(writer, http.StatusCreated, view)
}

func (h *Handler) findTicket(writer http.ResponseWriter, request *http.Request) {
	view := h.service.Find(request.Context(), request.PathValue("id"))
	if view.State == aftercare.StateError {
		status := http.StatusInternalServerError
		if view.Error != nil && view.Error.Code == "ticket_not_found" {
			status = http.StatusNotFound
		}
		writeJSON(writer, status, view)
		return
	}
	writeJSON(writer, http.StatusOK, view)
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
