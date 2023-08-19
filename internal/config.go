package internal

import "github.com/caarlos0/env/v7"

type Config struct {
	ServerAddr      string   `env:"SERVER_ADDR" envDefault:":3000"`
	Currencies      []string `env:"CURRENCIES" envDefault:"BRL,EUR,USD,INR"`
	CacheExpiration int64    `env:"CACHE_EXPIRATION" envDefault:"3600"`
	Redis           struct {
		Addr string `env:"REDIS_ADDR" envDefault:"127.0.0.1:6379"`
		User string `env:"REDIS_USER" envDefault:""`
		Pass string `env:"REDIS_PASS" envDefault:""`
	}
	API struct {
		URL   string `env:"API_URL" envDefault:"https://api.currencyapi.com/v3/latest"`
		Token string `env:"API_TOKEN" envDefault:"aszpkb7WFWtjBFxj9JHcorObU2vKTjaOiFCqmnAI"`
	}
}

// LoadConfig loads the config from the environment
func LoadConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
