package convert_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.come/edmarfelipe/currency-service/internal"
	"github.come/edmarfelipe/currency-service/internal/cache"
	"github.come/edmarfelipe/currency-service/internal/httpserver"
	"github.come/edmarfelipe/currency-service/internal/service"
	"github.come/edmarfelipe/currency-service/usecase/convert"
)

type test struct {
	name      string
	url       string
	expectErr error
}

func TestConvertController(t *testing.T) {
	ctrl := gomock.NewController(t)
	cacheService := cache.NewMockCache(ctrl)
	currencyService := service.NewMockCurrencyService(ctrl)

	cacheExpiration := int64(time.Minute * 60)
	cacheService.EXPECT().Get(gomock.Any(), "convert:BRL", cacheExpiration).Return(nil, nil).AnyTimes()
	cacheService.EXPECT().Set(gomock.Any(), "convert:BRL", gomock.Any(), cacheExpiration).AnyTimes()
	cacheService.EXPECT().Get(gomock.Any(), "convert:USD", cacheExpiration).Return(nil, nil).AnyTimes()
	cacheService.EXPECT().Set(gomock.Any(), "convert:USD", gomock.Any(), cacheExpiration).AnyTimes()

	currencyService.EXPECT().GetRate(gomock.Any(), "BRL").
		Return([]service.SymbolValue{
			{Code: "EUR", Value: 0.18168},
			{Code: "INR", Value: 16.230682},
			{Code: "USD", Value: 0.197499},
		}, nil).
		AnyTimes()

	currencyService.EXPECT().GetRate(gomock.Any(), "USD").
		Return([]service.SymbolValue{
			{Code: "EUR", Value: 0.919901},
			{Code: "INR", Value: 82.180893},
			{Code: "BRL", Value: 5.063305},
		}, nil).
		AnyTimes()

	currencyService.EXPECT().GetRate(gomock.Any(), "INR").
		Return([]service.SymbolValue{}, service.ErrFailedToConnect).
		AnyTimes()

	server := httpserver.New(&internal.Container{
		Config: &internal.Config{
			Currencies:      []string{"BRL", "EUR", "INR", "USD"},
			CacheExpiration: cacheExpiration,
		},
		CurrencyService: currencyService,
	}).TestServer(t)

	t.Run("Should validate invalid parameters", func(t *testing.T) {
		ts := []test{
			{
				"Should return 400 when value is invalid",
				server.URL + "/api/convert/USD/UNDEFINED",
				convert.ErrInvalidInputValue,
			},
			{
				"Should return 400 when currency is invalid",
				server.URL + "/api/convert/US/99",
				convert.ErrInvalidInputValue,
			},
			{
				"Should return 400 when currency is not supported",
				server.URL + "/api/convert/CAD/99",
				convert.ErrInvalidInputValue,
			},
		}

		for _, tc := range ts {
			t.Run(tc.name, func(t *testing.T) {
				resp, err := http.Get(tc.url)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			})
		}
	})

	t.Run("Should calculate correct when currency is BRL", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/convert/BRL/10")
		result := parseResult(resp)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 3, len(result))
		assert.Equal(t, "EUR", result[0].Currency)
		assert.Equal(t, "1.8168", result[0].Value)
		assert.Equal(t, "INR", result[1].Currency)
		assert.Equal(t, "162.3068", result[1].Value)
		assert.Equal(t, "USD", result[2].Currency)
		assert.Equal(t, "1.9750", result[2].Value)
	})

	t.Run("Should calculate correct when currency is USD", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/convert/USD/10")
		result := parseResult(resp)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 3, len(result))
		assert.Equal(t, 3, len(result))
		assert.Equal(t, "EUR", result[0].Currency)
		assert.Equal(t, "9.1990", result[0].Value)
		assert.Equal(t, "INR", result[1].Currency)
		assert.Equal(t, "821.8089", result[1].Value)
		assert.Equal(t, "BRL", result[2].Currency)
		assert.Equal(t, "50.6330", result[2].Value)
	})

	t.Run("Should handle correctly when an error occur on the use case", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/convert/INR/10")
		bytes, _ := io.ReadAll(resp.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "", string(bytes))
	})
}

func parseResult(resp *http.Response) []convert.CurrencyValue {
	var result []convert.CurrencyValue
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return result
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result
	}
	return result
}
