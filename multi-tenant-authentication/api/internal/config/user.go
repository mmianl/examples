package config

var globalUser UserConfig

type UserConfig struct {
}

// GlobalUser returns the globalUser config
func GlobalUser() *UserConfig {
	return &globalUser
}
