package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"tunnel-manager/models"
)

// Store manages application state with JSON file persistence
type Store struct {
	mu      sync.RWMutex
	filePath string
	config  models.Config
}

// NewStore creates a new Store with the given file path
func NewStore(filePath string) *Store {
	s := &Store{
		filePath: filePath,
		config: models.Config{
			PreferredCNAME: "cf.090227.xyz",
			AdminUsername:  "admin",
		},
	}
	s.load()
	return s
}

func (s *Store) load() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		// First run: generate random password, save config, print credentials
		password := generateRandomPassword(12)
		s.config.AdminPasswordHash = hashPassword(password)
		s.save()
		log.Printf("========================================")
		log.Printf("  首次启动，已生成管理员账户：")
		log.Printf("  用户名: %s", s.config.AdminUsername)
		log.Printf("  密  码: %s", password)
		log.Printf("  请登录后立即修改密码！")
		log.Printf("========================================")
		return
	}
	json.Unmarshal(data, &s.config)
	if s.config.PreferredCNAME == "" {
		s.config.PreferredCNAME = "cf.090227.xyz"
	}
	if s.config.AdminUsername == "" {
		s.config.AdminUsername = "admin"
	}
	if s.config.AdminPasswordHash == "" {
		password := generateRandomPassword(12)
		s.config.AdminPasswordHash = hashPassword(password)
		s.save()
		log.Printf("========================================")
		log.Printf("  密码为空，已自动生成：")
		log.Printf("  用户名: %s", s.config.AdminUsername)
		log.Printf("  密  码: %s", password)
		log.Printf("========================================")
	}
}

func (s *Store) save() {
	data, _ := json.MarshalIndent(s.config, "", "  ")
	os.WriteFile(s.filePath, data, 0644)
}

// GetConfig returns the current configuration
func (s *Store) GetConfig() models.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// SetTunnelID sets the active tunnel ID
func (s *Store) SetTunnelID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.TunnelID = id
	s.save()
}

// SetServiceURL sets the forwarding service URL
func (s *Store) SetServiceURL(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.ServiceURL = url
	s.save()
}

// SetPreferredCNAME sets the global preferred CNAME
func (s *Store) SetPreferredCNAME(cname string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.PreferredCNAME = cname
	s.save()
}

// GetAdminCredentials returns admin username and password hash
func (s *Store) GetAdminCredentials() (string, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.AdminUsername, s.config.AdminPasswordHash
}

// SetAdminCredentials sets admin username and password hash
func (s *Store) SetAdminCredentials(username, passwordHash string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.AdminUsername = username
	s.config.AdminPasswordHash = passwordHash
	s.save()
}

// ValidatePassword checks a plaintext password against the stored hash
func (s *Store) ValidatePassword(password, hash string) bool {
	return hashPassword(password) == hash
}

// HashPassword returns the SHA-256 hex digest of a password
func HashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", h)
}

// hashPassword is an internal alias for convenience
func hashPassword(password string) string {
	return HashPassword(password)
}

func generateRandomPassword(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)[:length]
}