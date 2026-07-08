package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"

	"tunnel-manager/models"
	"tunnel-manager/store"
)

// AdminHandler handles admin authentication and account management
type AdminHandler struct {
	store    *store.Store
	sessions map[string]bool
	mu       sync.RWMutex
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(st *store.Store) *AdminHandler {
	return &AdminHandler{
		store:    st,
		sessions: make(map[string]bool),
	}
}

// Login handles POST /api/admin/login
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		return
	}

	username, passwordHash := h.store.GetAdminCredentials()
	if !h.store.ValidatePassword(req.Password, passwordHash) || req.Username != username {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	token := h.generateToken()
	h.mu.Lock()
	h.sessions[token] = true
	h.mu.Unlock()

	writeJSON(w, http.StatusOK, models.LoginResponse{
		Token:    token,
		Username: username,
	})
}

// Status handles GET /api/admin/status — checks if a token is valid
func (h *AdminHandler) Status(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Auth-Token")
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"authenticated": "false"})
		return
	}

	h.mu.RLock()
	valid := h.sessions[token]
	h.mu.RUnlock()

	if !valid {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"authenticated": "false"})
		return
	}

	username, _ := h.store.GetAdminCredentials()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"username":      username,
	})
}

// ChangePassword handles PUT /api/admin/password
func (h *AdminHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req models.ChangePasswordRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	username, passwordHash := h.store.GetAdminCredentials()
	if !h.store.ValidatePassword(req.CurrentPassword, passwordHash) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "current password is incorrect"})
		return
	}

	if req.NewPassword == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "new password cannot be empty"})
		return
	}

	if len(req.NewPassword) < 6 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "new password must be at least 6 characters"})
		return
	}

	h.store.SetAdminCredentials(username, store.HashPassword(req.NewPassword))
	writeJSON(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
}

// ChangeUsername handles PUT /api/admin/username
func (h *AdminHandler) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	var req models.ChangeUsernameRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	_, passwordHash := h.store.GetAdminCredentials()
	if !h.store.ValidatePassword(req.CurrentPassword, passwordHash) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "password is incorrect"})
		return
	}

	if req.NewUsername == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "new username cannot be empty"})
		return
	}

	h.store.SetAdminCredentials(req.NewUsername, passwordHash)
	writeJSON(w, http.StatusOK, map[string]string{"message": "username updated successfully"})
}

// ValidateToken checks if a session token is valid (for middleware)
func (h *AdminHandler) ValidateToken(token string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.sessions[token]
}

// Logout handles POST /api/admin/logout
func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Auth-Token")
	if token != "" {
		h.mu.Lock()
		delete(h.sessions, token)
		h.mu.Unlock()
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *AdminHandler) generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}