package handlers

import (
	"io"
	"net/http"
	"strings"

	"tunnel-manager/models"
	"tunnel-manager/services"
	"tunnel-manager/store"
)

// TelegramHandler handles Telegram bot endpoints
type TelegramHandler struct {
	store *store.Store
	bot   *services.TelegramBot
}

// NewTelegramHandler creates a new TelegramHandler
func NewTelegramHandler(st *store.Store, bot *services.TelegramBot) *TelegramHandler {
	return &TelegramHandler{store: st, bot: bot}
}

// GetSettings returns the current bot settings (token masked)
func (h *TelegramHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	cfg := h.store.GetConfig()
	resp := models.TelegramSettingsResponse{
		Enabled:      cfg.TGBotEnabled,
		BotTokenSet:  cfg.TGBotToken != "",
		BotTokenHint: maskToken(cfg.TGBotToken),
		AdminTGIDs:   cfg.TGAdminIDs,
		Mode:         cfg.TGMode,
		WebhookURL:   cfg.TGWebhookURL,
		ApiEndpoint:  cfg.TGApiEndpoint,
	}
	writeJSON(w, http.StatusOK, resp)
}

// SaveSettings saves bot settings and restarts the bot
func (h *TelegramHandler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	var req models.TelegramSettingsRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// Validate
	if req.Enabled {
		// Check token: must be provided now or already stored
		cfg := h.store.GetConfig()
		if req.BotToken == "" && cfg.TGBotToken == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "启用 Bot 需要先设置 Token"})
			return
		}
		if req.AdminTGIDs == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "管理员 TG ID 不能为空"})
			return
		}
		if req.Mode == "webhook" && req.WebhookURL == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Webhook 模式需要设置公网访问地址"})
			return
		}
	}

	// Stop the bot before saving
	h.bot.Stop()

	// Save settings (empty token means keep existing)
	mode := req.Mode
	if mode == "" {
		mode = "polling"
	}
	apiEndpoint := req.ApiEndpoint
	if apiEndpoint == "" {
		cfg := h.store.GetConfig()
		apiEndpoint = cfg.TGApiEndpoint
	}
	if apiEndpoint == "" {
		apiEndpoint = "https://api.telegram.org"
	}

	h.store.SetTelegramSettings(req.Enabled, req.BotToken, req.AdminTGIDs, mode, req.WebhookURL, apiEndpoint)

	// Restart if enabled
	var startErr error
	if req.Enabled {
		startErr = h.bot.Start()
	}

	resp := map[string]interface{}{
		"status":  "ok",
		"running": h.bot.Status().Running,
	}
	if startErr != nil {
		resp["error"] = startErr.Error()
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetStatus returns bot running status
func (h *TelegramHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.bot.Status())
}

// SendTest sends a test message to all admin IDs
func (h *TelegramHandler) SendTest(w http.ResponseWriter, r *http.Request) {
	if err := h.bot.SendTestMessage(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "测试消息已发送"})
}

// Webhook handles incoming Telegram webhook updates (no auth middleware)
func (h *TelegramHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	cfg := h.store.GetConfig()
	if !cfg.TGBotEnabled || cfg.TGMode != "webhook" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "bot not in webhook mode"})
		return
	}

	secret := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if !h.bot.VerifyWebhookSecret(secret) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "invalid secret token"})
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "read body failed"})
		return
	}

	h.bot.HandleWebhookUpdate(body)
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// maskToken returns a masked version of the bot token
func maskToken(token string) string {
	if token == "" {
		return ""
	}
	parts := strings.SplitN(token, ":", 2)
	if len(parts) != 2 || len(parts[1]) <= 4 {
		return parts[0] + ":****"
	}
	return parts[0] + ":" + strings.Repeat("*", len(parts[1])-4) + parts[1][len(parts[1])-4:]
}
