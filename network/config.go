package network

type Config struct {
	base_url string
	timeout  int
}

func (c *Config) ValidateStatus(status int) bool {
	return true
}
