package models

import "encoding/json"

// Config represents the application configuration state
type Config struct {
	TunnelID           string `json:"tunnel_id"`
	ServiceURL         string `json:"service_url"`
	PreferredCNAME     string `json:"preferred_cname"`
	AdminUsername      string `json:"admin_username"`
	AdminPasswordHash  string `json:"admin_password_hash"`
}

// Tunnel represents a Cloudflare Tunnel
type Tunnel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Zone represents a Cloudflare Zone
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DNSRecord represents a DNS record
type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
}

// TunnelConfigResponse represents the CF API response for tunnel config
type TunnelConfigResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Config struct {
			Ingress []IngressRule `json:"ingress"`
		} `json:"config"`
	} `json:"result"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// IngressRule represents a tunnel ingress rule
type IngressRule struct {
	Hostname string `json:"hostname,omitempty"`
	Service  string `json:"service"`
}

// CFAPIResponse is a generic Cloudflare API response wrapper
type CFAPIResponse struct {
	Success  bool            `json:"success"`
	Result   json.RawMessage `json:"result"`
	Errors   []CFError       `json:"errors"`
	Messages []string        `json:"messages"`
}

// CFError represents a Cloudflare API error
type CFError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BindRequest is the request body for domain binding
type BindRequest struct {
	MainDomain string `json:"main_domain"`
	AuxDomain  string `json:"aux_domain"`
}

// FallbackRequest is the request body for setting fallback origin
type FallbackRequest struct {
	Domain string `json:"domain"`
}

// CustomHostname represents a Cloudflare SaaS custom hostname
type CustomHostname struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
}

// SetValueRequest is a generic request for setting a single value
type SetValueRequest struct {
	Value string `json:"value"`
}

// LoginRequest is the request body for admin login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse is the response body for admin login
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

// ChangePasswordRequest is the request body for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangeUsernameRequest is the request body for changing username
type ChangeUsernameRequest struct {
	CurrentPassword string `json:"current_password"`
	NewUsername     string `json:"new_username"`
}