package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"tunnel-manager/handlers"
	"tunnel-manager/services"
	"tunnel-manager/store"
)

func main() {
	apiToken := os.Getenv("CF_API_TOKEN")
	accountID := os.Getenv("CF_ACCOUNT_ID")
	apiKey := os.Getenv("API_KEY")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	storePath := os.Getenv("STORE_PATH")
	if storePath == "" {
		storePath = "data/config.json"
	}

	if apiToken == "" || accountID == "" {
		log.Fatal("CF_API_TOKEN and CF_ACCOUNT_ID environment variables are required")
	}

	// Initialize dependencies
	st := store.NewStore(storePath)
	cf := services.NewCloudflareClient(apiToken, accountID)

	// Initialize handlers
	configHandler := handlers.NewConfigHandler(st)
	tunnelHandler := handlers.NewTunnelHandler(cf)
	domainHandler := handlers.NewDomainHandler(cf, st)
	adminHandler := handlers.NewAdminHandler(st)

	mw := &handlers.Middleware{
		APIKey:       apiKey,
		AdminHandler: adminHandler,
	}

	// Setup router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(mw.CORS)

	r.Route("/api", func(r chi.Router) {
		// Admin endpoints (no auth required)
		r.Post("/admin/login", adminHandler.Login)
		r.Post("/admin/logout", adminHandler.Logout)
		r.Get("/admin/status", adminHandler.Status)

		// Protected admin endpoints (behind auth)
		r.Put("/admin/password", mw.Auth(adminHandler.ChangePassword))
		r.Put("/admin/username", mw.Auth(adminHandler.ChangeUsername))

		// Config endpoints
		r.Get("/config", mw.Auth(configHandler.GetConfig))
		r.Post("/config/tunnel", mw.Auth(configHandler.SetTunnelID))
		r.Post("/config/service", mw.Auth(configHandler.SetServiceURL))
		r.Post("/config/preferred-cname", mw.Auth(configHandler.SetPreferredCNAME))

		// Tunnel endpoints
		r.Get("/tunnels", mw.Auth(tunnelHandler.ListTunnels))
		r.Get("/zones", mw.Auth(tunnelHandler.ListZones))

		// Domain endpoints
		r.Post("/domain/bind", mw.Auth(domainHandler.BindDomain))
		r.Post("/domain/fallback", mw.Auth(domainHandler.SetFallbackOrigin))

		// Health check (no auth)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})
	})

	// Serve frontend static files (SPA fallback)
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "frontend/dist"
	}
	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		fs := http.FileServer(http.Dir(staticDir))
		r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			// Try to serve the file directly
			path := filepath.Join(staticDir, req.URL.Path)
			if _, err := os.Stat(path); os.IsNotExist(err) || strings.HasSuffix(req.URL.Path, "/") {
				// SPA fallback: serve index.html for missing routes
				http.ServeFile(w, req, filepath.Join(staticDir, "index.html"))
				return
			}
			fs.ServeHTTP(w, req)
		})
		log.Printf("Serving static files from %s", staticDir)
	}

	addr := ":" + port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}