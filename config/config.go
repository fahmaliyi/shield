package config

import "time"

const DefaultSessionTTL = 30 * time.Minute

type Config struct {
	SessionTTL time.Duration
}

func DefaultConfig() Config {
	return Config{
		SessionTTL: DefaultSessionTTL,
	}
}
