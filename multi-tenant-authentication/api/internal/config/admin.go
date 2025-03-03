package config

var globalAdmin AdminConfig

type AdminConfig struct {
}

// GlobalAdmin returns the globalAdmin config
func GlobalAdmin() *AdminConfig {
	return &globalAdmin
}
