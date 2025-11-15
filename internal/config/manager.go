package config

import (
	"sync"

	"github.com/happytaoer/prompt-security/internal/db"
)

// Manager manages configuration with dynamic reload support
type Manager struct {
	config   Config
	mu       sync.RWMutex
	onChange []func(Config) // Callbacks to notify when config changes
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	return &Manager{
		config:   cfg,
		onChange: make([]func(Config), 0),
	}, nil
}

// Get returns a copy of the current configuration
func (m *Manager) Get() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// Update updates the configuration and notifies all listeners
func (m *Manager) Update(cfg Config) error {
	// Save to database first
	if err := db.SaveConfig(cfg); err != nil {
		return err
	}

	// Update in-memory config
	m.mu.Lock()
	m.config = cfg
	callbacks := m.onChange
	m.mu.Unlock()

	// Notify all listeners
	for _, callback := range callbacks {
		callback(cfg)
	}

	return nil
}

// OnChange registers a callback to be called when configuration changes
func (m *Manager) OnChange(callback func(Config)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChange = append(m.onChange, callback)
}

// Reload reloads configuration from database
func (m *Manager) Reload() error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.config = cfg
	callbacks := m.onChange
	m.mu.Unlock()

	// Notify all listeners
	for _, callback := range callbacks {
		callback(cfg)
	}

	return nil
}
