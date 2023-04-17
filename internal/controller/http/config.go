package http

type Config struct {
	BindAddr string `toml:"bind_add"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}
