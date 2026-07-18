package handlers

import (
	"net/http"

	"tunnel-manager/models"
	"tunnel-manager/store"
)

// ConfigHandler handles configuration-related requests
type ConfigHandler struct {
	store *store.Store
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(s *store.Store) *ConfigHandler {
	return &ConfigHandler{store: s}
}

// GetConfig returns the current configuration (sanitized)
func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := h.store.GetConfig()
	writeJSON(w, http.StatusOK, map[string]string{
		"tunnel_id":       cfg.TunnelID,
		"service_url":     cfg.ServiceURL,
		"preferred_cname": cfg.PreferredCNAME,
	})
}

// SetTunnelID sets the active tunnel ID
func (h *ConfigHandler) SetTunnelID(w http.ResponseWriter, r *http.Request) {
	var req models.SetValueRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Value == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "value is required"})
		return
	}
	h.store.SetTunnelID(req.Value)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "tunnel_id": req.Value})
}

// SetServiceURL sets the forwarding service URL
func (h *ConfigHandler) SetServiceURL(w http.ResponseWriter, r *http.Request) {
	var req models.SetValueRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Value == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "value is required"})
		return
	}
	h.store.SetServiceURL(req.Value)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service_url": req.Value})
}

// SetPreferredCNAME sets the global preferred CNAME
func (h *ConfigHandler) SetPreferredCNAME(w http.ResponseWriter, r *http.Request) {
	var req models.SetValueRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Value == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "value is required"})
		return
	}
	h.store.SetPreferredCNAME(req.Value)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "preferred_cname": req.Value})
}