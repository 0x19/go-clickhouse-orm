package chorm

import "fmt"

type Config struct {
	Host     string `json:"host" yaml:"host" env:"GCHM_HOST"`
	Port     int    `json:"port" yaml:"port" env:"GCHM_PORT"`
	Username string `json:"username" yaml:"username" env:"GCHM_USERNAME"`
	Password string `json:"password" yaml:"password" env:"GCHM_PASSWORD"`
	Database string `json:"database" yaml:"database" env:"GCHM_DATABASE"`
	SSLMode  string `json:"ssl_mode" yaml:"ssl_mode" env:"GCHM_SSL_MODE"`
	Insecure bool   `json:"insecure" yaml:"insecure" env:"GCHM_INSECURE"`
}

// Validate validates the configuration provided.
// TOOD: In the future, this should be replaced with a more robust validation
func (c *Config) Validate() error {
	if c == nil {
		return ErrNoConfigProvided
	}

	if c.Host == "" {
		return ErrNoHostProvided
	}

	if c.Port == 0 {
		return ErrNoPortProvided
	}

	if c.Username == "" {
		return ErrNoUsernameProvided
	}

	if c.Password == "" {
		return ErrNoPasswordProvided
	}

	if c.Database == "" {
		return ErrNoDatabaseProvided
	}

	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
		c.SSLMode,
	)
}
