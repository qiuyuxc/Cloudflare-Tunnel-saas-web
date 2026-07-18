package handlers

import (
	"fmt"
	"net/http"

	"tunnel-manager/models"
	"tunnel-manager/services"
)

// DomainHandler handles domain binding and fallback origin requests
type DomainHandler struct {
	svc *services.DomainService
}

// NewDomainHandler creates a new DomainHandler
func NewDomainHandler(svc *services.DomainService) *DomainHandler {
	return &DomainHandler{svc: svc}
}

// BindDomain performs the full domain binding flow
func (h *DomainHandler) BindDomain(w http.ResponseWriter, r *http.Request) {
	var req models.BindRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.MainDomain == "" || req.AuxDomain == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "main_domain and aux_domain are required"})
		return
	}

	_, err := h.svc.BindDomain(req.MainDomain, req.AuxDomain)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":      "ok",
		"message":     fmt.Sprintf("Domain binding complete! Access: https://%s", req.MainDomain),
		"main_domain": req.MainDomain,
		"aux_domain":  req.AuxDomain,
	})
}

// SetFallbackOrigin sets the fallback origin for custom hostnames
func (h *DomainHandler) SetFallbackOrigin(w http.ResponseWriter, r *http.Request) {
	var req models.FallbackRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Domain == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "domain is required"})
		return
	}

	if err := h.svc.SetFallbackOrigin(req.Domain); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": fmt.Sprintf("Fallback origin set to %s", req.Domain),
	})
}
