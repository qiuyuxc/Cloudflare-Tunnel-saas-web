package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tunnel-manager/models"
)

// CloudflareClient wraps Cloudflare API calls
type CloudflareClient struct {
	apiToken  string
	accountID string
	baseURL   string
	httpClient *http.Client
}

// NewCloudflareClient creates a new Cloudflare API client
func NewCloudflareClient(apiToken, accountID string) *CloudflareClient {
	return &CloudflareClient{
		apiToken:  apiToken,
		accountID: accountID,
		baseURL:   "https://api.cloudflare.com/client/v4",
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *CloudflareClient) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *CloudflareClient) do(req *http.Request, target interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}

	var apiResp models.CFAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("parse response failed: %w", err)
	}

	if !apiResp.Success {
		errMsg := "unknown error"
		if len(apiResp.Errors) > 0 {
			errMsg = apiResp.Errors[0].Message
		}
		return fmt.Errorf("cloudflare API error: %s", errMsg)
	}

	if target != nil {
		return json.Unmarshal(apiResp.Result, target)
	}

	return nil
}

// ListTunnels lists all Cloudflare tunnels
func (c *CloudflareClient) ListTunnels() ([]models.Tunnel, error) {
	path := fmt.Sprintf("/accounts/%s/cfd_tunnel?is_deleted=false", c.accountID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var tunnels []models.Tunnel
	if err := c.do(req, &tunnels); err != nil {
		return nil, err
	}
	return tunnels, nil
}

// GetTunnelConfig fetches the current tunnel configuration
func (c *CloudflareClient) GetTunnelConfig(tunnelID string) (*models.TunnelConfigResponse, error) {
	path := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s/configurations", c.accountID, tunnelID)
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var cfg models.TunnelConfigResponse
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if err := json.Unmarshal(body, &cfg); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	if !cfg.Success {
		errMsg := "unknown error"
		if len(cfg.Errors) > 0 {
			errMsg = cfg.Errors[0].Message
		}
		return nil, fmt.Errorf("Cloudflare API: %s", errMsg)
	}

	return &cfg, nil
}

// UpdateTunnelConfig updates the tunnel configuration
func (c *CloudflareClient) UpdateTunnelConfig(tunnelID string, config interface{}) error {
	path := fmt.Sprintf("/accounts/%s/cfd_tunnel/%s/configurations", c.accountID, tunnelID)
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req, err := c.newRequest("PUT", path, strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// GetZoneIDByHostname finds the zone ID for a given hostname
func (c *CloudflareClient) GetZoneIDByHostname(hostname string) (string, error) {
	hostname = strings.TrimSpace(strings.ToLower(strings.TrimSuffix(hostname, ".")))

	req, err := c.newRequest("GET", "/zones?status=active&per_page=1000", nil)
	if err != nil {
		return "", err
	}

	var zones []models.Zone
	if err := c.do(req, &zones); err != nil {
		return "", err
	}

	var bestMatch *models.Zone
	for _, z := range zones {
		zoneName := strings.TrimSpace(strings.ToLower(z.Name))
		if hostname == zoneName || strings.HasSuffix(hostname, "."+zoneName) {
			if bestMatch == nil || len(zoneName) > len(bestMatch.Name) {
				zCopy := z
				bestMatch = &zCopy
			}
		}
	}

	if bestMatch == nil {
		return "", fmt.Errorf("no zone found for hostname: %s", hostname)
	}
	return bestMatch.ID, nil
}

// UpsertDNSRecord creates or updates a DNS record
func (c *CloudflareClient) UpsertDNSRecord(zoneID, name, recordType, content string, proxied bool) error {
	// List existing records
	listURL := fmt.Sprintf("/zones/%s/dns_records?name=%s&type=%s", zoneID, url.QueryEscape(name), recordType)
	req, err := c.newRequest("GET", listURL, nil)
	if err != nil {
		return err
	}

	var records []models.DNSRecord
	if err := c.do(req, &records); err != nil {
		// If listing fails, try to create directly
		records = nil
	}

	payload := models.DNSRecord{
		Name:    name,
		Type:    recordType,
		Content: content,
		Proxied: proxied,
	}
	body, _ := json.Marshal(payload)

	if len(records) > 0 {
		// Update existing record
		updateURL := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, records[0].ID)
		req, err = c.newRequest("PUT", updateURL, strings.NewReader(string(body)))
	} else {
		// Create new record
		createURL := fmt.Sprintf("/zones/%s/dns_records", zoneID)
		req, err = c.newRequest("POST", createURL, strings.NewReader(string(body)))
	}

	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// SetFallbackOrigin sets the fallback origin for custom hostnames
func (c *CloudflareClient) SetFallbackOrigin(zoneID, origin string) error {
	path := fmt.Sprintf("/zones/%s/custom_hostnames/fallback_origin", zoneID)
	payload := map[string]string{"origin": origin}
	body, _ := json.Marshal(payload)

	req, err := c.newRequest("PUT", path, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// CreateCustomHostname creates a SaaS custom hostname
func (c *CloudflareClient) CreateCustomHostname(zoneID, hostname, originServer string) error {
	path := fmt.Sprintf("/zones/%s/custom_hostnames", zoneID)
	payload := map[string]interface{}{
		"hostname":             hostname,
		"custom_origin_server": originServer,
		"ssl": map[string]interface{}{
			"method": "http",
			"type":   "dv",
		},
	}
	body, _ := json.Marshal(payload)

	req, err := c.newRequest("POST", path, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ListZones lists all active zones
func (c *CloudflareClient) ListZones() ([]models.Zone, error) {
	req, err := c.newRequest("GET", "/zones?status=active&per_page=1000", nil)
	if err != nil {
		return nil, err
	}

	var zones []models.Zone
	if err := c.do(req, &zones); err != nil {
		return nil, err
	}
	return zones, nil
}