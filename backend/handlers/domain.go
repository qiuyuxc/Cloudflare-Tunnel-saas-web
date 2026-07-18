package handlers

import (
	"fmt"
	"net/http"

	"tunnel-manager/models"
	"tunnel-manager/services"
	"tunnel-manager/store"
)

// DomainHandler handles domain binding and fallback origin requests
type DomainHandler struct {
	cf    *services.CloudflareClient
	store *store.Store
}

// NewDomainHandler creates a new DomainHandler
func NewDomainHandler(cf *services.CloudflareClient, s *store.Store) *DomainHandler {
	return &DomainHandler{cf: cf, store: s}
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

	cfg := h.store.GetConfig()
	if cfg.TunnelID == "" || cfg.ServiceURL == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "tunnel_id and service_url must be configured first. Use /api/config/tunnel and /api/config/service endpoints.",
		})
		return
	}

	// Get zone IDs
	mainZoneID, err := h.cf.GetZoneIDByHostname(req.MainDomain)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("main domain: %s", err.Error())})
		return
	}

	auxZoneID, err := h.cf.GetZoneIDByHostname(req.AuxDomain)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("aux domain: %s", err.Error())})
		return
	}

	// Update tunnel ingress rules
	tunnelCfg, err := h.cf.GetTunnelConfig(cfg.TunnelID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("get tunnel config: %s", err.Error())})
		return
	}

	newRules := []models.IngressRule{
		{Hostname: req.MainDomain, Service: cfg.ServiceURL},
		{Hostname: req.AuxDomain, Service: cfg.ServiceURL},
	}

	// Deduplicate: remove existing rules for these hostnames, then insert updated ones
	ingress := tunnelCfg.Result.Config.Ingress
	filtered := make([]models.IngressRule, 0, len(ingress))
	for _, rule := range ingress {
		if rule.Hostname != req.MainDomain && rule.Hostname != req.AuxDomain {
			filtered = append(filtered, rule)
		}
	}
	if len(filtered) > 0 {
		lastRule := filtered[len(filtered)-1]
		filtered = filtered[:len(filtered)-1]
		filtered = append(filtered, newRules...)
		filtered = append(filtered, lastRule)
	} else {
		filtered = append(newRules, models.IngressRule{Service: "http_status:404"})
	}
	tunnelCfg.Result.Config.Ingress = filtered

	updatePayload := map[string]interface{}{
		"config": tunnelCfg.Result.Config,
	}

	if err := h.cf.UpdateTunnelConfig(cfg.TunnelID, updatePayload); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("update tunnel config: %s", err.Error())})
		return
	}

	// Upsert DNS records
	tunnelCNAME := fmt.Sprintf("%s.cfargotunnel.com", cfg.TunnelID)
	if err := h.cf.UpsertDNSRecord(auxZoneID, req.AuxDomain, "CNAME", tunnelCNAME, true); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("upsert aux DNS: %s", err.Error())})
		return
	}

	if err := h.cf.UpsertDNSRecord(mainZoneID, req.MainDomain, "CNAME", cfg.PreferredCNAME, false); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("upsert main DNS: %s", err.Error())})
		return
	}

	// Upsert SaaS custom hostname
	if err := h.cf.UpsertCustomHostname(auxZoneID, req.MainDomain, req.AuxDomain); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("upsert custom hostname: %s", err.Error())})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":   "ok",
		"message":  fmt.Sprintf("Domain binding complete! Access: https://%s", req.MainDomain),
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

	zoneID, err := h.cf.GetZoneIDByHostname(req.Domain)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := h.cf.SetFallbackOrigin(zoneID, req.Domain); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"message": fmt.Sprintf("Fallback origin set to %s", req.Domain),
	})
}