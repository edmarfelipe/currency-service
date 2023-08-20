package internal

import (
	"log/slog"

	"github.com/rueian/rueidis"
	"github.come/edmarfelipe/currency-service/internal/cache"
	"github.come/edmarfelipe/currency-service/internal/service"
)

// Container is a container for the application. that holds all the dependencies
type Container struct {
	redisClient     rueidis.Client
	Config          *Config
	CurrencyService service.CurrencyService
}

// NewContainer creates a new container with all the dependencies
func NewContainer(cfg *Config) (*Container, error) {
	redisClient, err := cache.Open(cfg.Redis.Addr, cfg.Redis.User, cfg.Redis.Pass)
	if err != nil {
		return nil, err
	}

	return &Container{
		Config:      cfg,
		redisClient: redisClient,
		CurrencyService: service.NewCurrencyCache(
			service.NewCurrencyService(cfg.API.URL, cfg.API.Token, cfg.Currencies),
			cache.New(redisClient),
			cfg.CacheExpiration,
		),
	}, nil
}

// Shutdown closes all the connections
func (c *Container) Shutdown() {
	slog.Info("Closing redis connections")
	c.redisClient.Close()
	slog.Info("Redis connections closed")
}
