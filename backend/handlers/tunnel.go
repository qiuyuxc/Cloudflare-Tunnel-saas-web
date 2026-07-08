package handlers

import (
	"net/http"

	"tunnel-manager/services"
)

// TunnelHandler handles tunnel-related requests
type TunnelHandler struct {
	cf *services.CloudflareClient
}

// NewTunnelHandler creates a new TunnelHandler
func NewTunnelHandler(cf *services.CloudflareClient) *TunnelHandler {
	return &TunnelHandler{cf: cf}
}

// ListTunnels returns all Cloudflare tunnels
func (h *TunnelHandler) ListTunnels(w http.ResponseWriter, r *http.Request) {
	tunnels, err := h.cf.ListTunnels()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, tunnels)
}

// ListZones returns all Cloudflare zones
func (h *TunnelHandler) ListZones(w http.ResponseWriter, r *http.Request) {
	zones, err := h.cf.ListZones()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, zones)
}