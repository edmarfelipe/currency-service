package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_currencyCache_GetRate(t *testing.T) {

	t.Run("Should return from cache when cache is not about to expire", func(t *testing.T) {})

	t.Run("Should return from service when cache is expired", func(t *testing.T) {})

	t.Run("Should update cache when cache is about to expire", func(t *testing.T) {})

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
