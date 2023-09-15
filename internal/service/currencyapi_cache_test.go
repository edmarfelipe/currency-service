package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.come/edmarfelipe/currency-service/internal/cache"
)

func Test_currencyCache_GetRate(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("Should return from cache when cache is not about to expire", func(t *testing.T) {
		cacheService := cache.NewMockCache(ctrl)

		srv := &currencyCache{
			cacheExpiration: (1 * time.Minute).Milliseconds(),
			cache:           cacheService,
		}

		cacheService.EXPECT().
			Get(gomock.Any(), "convert:BRL", gomock.Any()).
			Return(&cache.Result{
				TTLRemaining: (1 * time.Minute).Milliseconds(),
				Value:        []byte(`[{ "code": "BRL", "value": 3.5}]`),
			}, nil).
			AnyTimes()

		values, err := srv.GetRate(context.Background(), "BRL")

		assert.Nil(t, err)
		assert.Len(t, values, 1)
		assert.Contains(t, values, SymbolValue{Code: "BRL", Value: 3.5})
	})

	t.Run("Should return from service when cache is expired", func(t *testing.T) {
		currencyService := NewMockCurrencyService(ctrl)
		cacheService := cache.NewMockCache(ctrl)

		srv := NewCurrencyCache(currencyService, cacheService, (1 * time.Minute).Milliseconds())

		cacheService.EXPECT().
			Get(gomock.Any(), "convert:EUR", gomock.Any()).Return(nil, nil).AnyTimes()

		cacheService.EXPECT().
			Set(gomock.Any(), "convert:EUR", gomock.Any(), (1 * time.Minute).Milliseconds()).Return(nil).AnyTimes()

		currencyService.EXPECT().GetRate(gomock.Any(), "EUR").
			Return([]SymbolValue{{Code: "EUR", Value: 0.1}}, nil).
			AnyTimes()

		result, err := srv.GetRate(context.Background(), "EUR")

		assert.Nil(t, err)
		assert.Equal(t, []SymbolValue{{Code: "EUR", Value: 0.1}}, result)
	})

	t.Run("Should update cache when cache is about to expire", func(t *testing.T) {
		currencyService := NewMockCurrencyService(ctrl)
		cacheService := cache.NewMockCache(ctrl)

		srv := NewCurrencyCache(currencyService, cacheService, (1 * time.Minute).Milliseconds())

		cacheService.EXPECT().
			Get(gomock.Any(), "convert:EUR", gomock.Any()).Return(nil, nil).AnyTimes()

		cacheService.EXPECT().
			Set(gomock.Any(), "convert:EUR", gomock.Any(), (1 * time.Minute).Milliseconds()).Return(nil).AnyTimes()

		currencyService.EXPECT().GetRate(gomock.Any(), "EUR").
			Return([]SymbolValue{{Code: "EUR", Value: 0.1}}, nil).
			AnyTimes()

		_, err := srv.GetRate(context.Background(), "EUR")

		assert.Nil(t, err)
	})

	t.Run("isAboutToExpire", func(t *testing.T) {
		srv := &currencyCache{
			cacheExpiration: (100 * time.Minute).Milliseconds(),
		}

		t.Run("Should be true when TTL is 5% or less to the expiration time", func(t *testing.T) {
			assert.True(t, srv.isAboutToExpire((time.Minute * 5).Milliseconds()))
		})

		t.Run("Should be false when TTL is more than 5% to the expiration time", func(t *testing.T) {
			assert.False(t, srv.isAboutToExpire((time.Minute * 6).Milliseconds()))
		})
	})
}
