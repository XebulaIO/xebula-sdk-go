package network

type Config struct {
	BaseURL string
	Timeout int
}

func (c *Config) ValidateStatus(status int) bool {
	return true
}
