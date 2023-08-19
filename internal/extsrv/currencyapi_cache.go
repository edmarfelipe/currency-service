package extsrv

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.come/edmarfelipe/currency-service/internal/cache"
)

var (
	cacheKey           = "convert:"
	cacheMinPercentage = 5
)

type currencyCache struct {
	cacheExpiration int64
	srv             CurrencyService
	cache           cache.Cache
}

func NewCurrencyCache(srv CurrencyService, cache cache.Cache, cacheExpiration int64) CurrencyService {
	return &currencyCache{
		cacheExpiration: cacheExpiration,
		srv:             srv,
		cache:           cache,
	}
}

func (c *currencyCache) isAboutToExpire(ttl int64) bool {
	value := (int64(cacheMinPercentage) * c.cacheExpiration) / 100
	slog.Info("Calculating if cache is about to expire", "current-ttl", ttl, "ttl-before-expire", value)
	return value >= ttl
}

func (c *currencyCache) getFromService(ctx context.Context, currency string) ([]SymbolValue, error) {
	slog.InfoContext(ctx, "Getting rate from service")
	rate, err := c.srv.GetRate(ctx, currency)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "Saving rate to cache")
	err = c.cache.Set(ctx, cacheKey+currency, rate, c.cacheExpiration)
	if err != nil {
		slog.ErrorContext(ctx, "error saving cache", "err", err)
	}

	return rate, nil
}

func (c *currencyCache) getFromCache(ctx context.Context, currency string) []SymbolValue {
	slog.InfoContext(ctx, "Getting rate from cache")
	result, err := c.cache.Get(ctx, cacheKey+currency, c.cacheExpiration)
	if err != nil {
		slog.ErrorContext(ctx, "error getting cache", "err", err)
	}
	if result == nil {
		slog.ErrorContext(ctx, "Cache not found", "currency", currency)
		return nil
	}
	// When the cache is about to expire, we start a goroutine to get the rate from the service
	if c.isAboutToExpire(result.TTLRemaining) {
		go func(ctx context.Context, currency string) {
			slog.InfoContext(ctx, "Getting rate from service async")
			_, err = c.getFromService(ctx, currency)
			if err != nil {
				slog.ErrorContext(ctx, "error getting rate from service", "err", err)
			}
		}(context.WithoutCancel(ctx), currency)
	}

	var value []SymbolValue
	err = json.Unmarshal(result.Value, &value)
	if err != nil {
		slog.ErrorContext(ctx, "error unmarshalling cache", "err", err)
		return nil
	}
	slog.InfoContext(ctx, "Returning rate from cache", "ttl-remaining", result.TTLRemaining)
	return value
}

func (c *currencyCache) GetRate(ctx context.Context, currency string) ([]SymbolValue, error) {
	slog.InfoContext(ctx, "Checking cache for rate", "currency", currency)
	value := c.getFromCache(ctx, currency)
	if value != nil {
		return value, nil
	}

	rate, err := c.getFromService(ctx, currency)
	if err != nil {
		return nil, err
	}

	return rate, nil
}
