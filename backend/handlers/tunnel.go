package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"tunnel-manager/models"
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

// GetTunnelDetail returns tunnel details including ingress rules
func (h *TunnelHandler) GetTunnelDetail(w http.ResponseWriter, r *http.Request) {
	tunnelID := chi.URLParam(r, "tunnelID")

	tunnels, err := h.cf.ListTunnels()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var tunnelName, tunnelStatus string
	for _, t := range tunnels {
		if t.ID == tunnelID {
			tunnelName = t.Name
			tunnelStatus = t.Status
			break
		}
	}
	if tunnelName == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "tunnel not found"})
		return
	}

	cfg, err := h.cf.GetTunnelConfig(tunnelID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":      tunnelID,
		"name":    tunnelName,
		"status":  tunnelStatus,
		"ingress": cfg.Result.Config.Ingress,
	})
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

// IngressRuleRequest is the request body for adding/updating an ingress rule
type IngressRuleRequest struct {
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
}

// AddIngressRule adds a new ingress rule to the tunnel
func (h *TunnelHandler) AddIngressRule(w http.ResponseWriter, r *http.Request) {
	tunnelID := chi.URLParam(r, "tunnelID")

	var req IngressRuleRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Hostname == "" || req.Service == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "hostname and service are required"})
		return
	}

	cfg, err := h.cf.GetTunnelConfig(tunnelID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("get tunnel config: %s", err.Error())})
		return
	}

	ingress := cfg.Result.Config.Ingress
	// Check duplicate
	for _, rule := range ingress {
		if rule.Hostname == req.Hostname {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "hostname already exists, use PUT to update"})
			return
		}
	}

	// Insert before catch-all (last rule)
	if len(ingress) > 0 {
		catchall := ingress[len(ingress)-1]
		ingress = ingress[:len(ingress)-1]
		ingress = append(ingress, models.IngressRule{Hostname: req.Hostname, Service: req.Service})
		ingress = append(ingress, catchall)
	} else {
		ingress = []models.IngressRule{
			{Hostname: req.Hostname, Service: req.Service},
			{Service: "http_status:404"},
		}
	}
	cfg.Result.Config.Ingress = ingress

	if err := h.cf.UpdateTunnelConfig(tunnelID, map[string]interface{}{"config": cfg.Result.Config}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("update tunnel config: %s", err.Error())})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "route added"})
}

// UpdateIngressRule updates an existing ingress rule by hostname
func (h *TunnelHandler) UpdateIngressRule(w http.ResponseWriter, r *http.Request) {
	tunnelID := chi.URLParam(r, "tunnelID")

	var req struct {
		OldHostname string `json:"old_hostname"`
		Hostname    string `json:"hostname"`
		Service     string `json:"service"`
	}
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.OldHostname == "" || req.Hostname == "" || req.Service == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "old_hostname, hostname and service are required"})
		return
	}

	cfg, err := h.cf.GetTunnelConfig(tunnelID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("get tunnel config: %s", err.Error())})
		return
	}

	found := false
	for i, rule := range cfg.Result.Config.Ingress {
		if rule.Hostname == req.OldHostname {
			cfg.Result.Config.Ingress[i] = models.IngressRule{Hostname: req.Hostname, Service: req.Service}
			found = true
			break
		}
	}
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "route not found"})
		return
	}

	if err := h.cf.UpdateTunnelConfig(tunnelID, map[string]interface{}{"config": cfg.Result.Config}); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("update tunnel config: %s", err.Error())})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "route updated"})
}