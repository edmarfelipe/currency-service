package service_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.come/edmarfelipe/currency-service/internal/service"
)

const (
	baseURL = "https://api.fake.com/v3/latest"
	token   = ""
)

func createMockAPI() func() {
	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		baseURL+"?base_currency=BRL&currencies=EUR,USD,INR",
		httpmock.NewStringResponder(200, `{
			"data": {
				"EUR": {
					"code": "EUR",
					"value": 0.178613
				},
				"INR": {
					"code": "INR",
					"value": 15.90579
				},
				"USD": {
					"code": "USD",
					"value": 0.193584
				}
			}
		}`),
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		baseURL+"?base_currency=EUR&currencies=BRL",
		httpmock.NewStringResponder(200, `{
			"data": {
				"BRL": {
					"code": "BRL",
					"value": 5.165707
				}
			}
		}`),
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		baseURL+"?base_currency=LLL&currencies=BRL",
		httpmock.NewStringResponder(422, ""),
	)
	return httpmock.DeactivateAndReset
}

func TestNewCurrencyService(t *testing.T) {
	clearMock := createMockAPI()

	t.Run("Should get rate for only 1 currency", func(t *testing.T) {
		srv := service.NewCurrencyService(baseURL, token, []string{"BRL"})

		result, err := srv.GetRate(context.Background(), "EUR")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "BRL", result[0].Code)
		assert.Equal(t, 5.165707, result[0].Value)
	})

	t.Run("Should get rate for 3 currencies", func(t *testing.T) {
		srv := service.NewCurrencyService(baseURL, token, []string{"EUR", "USD", "INR"})

		result, err := srv.GetRate(context.Background(), "BRL")

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Contains(t, result, service.SymbolValue{Code: "EUR", Value: 0.178613})
		assert.Contains(t, result, service.SymbolValue{Code: "INR", Value: 15.90579})
		assert.Contains(t, result, service.SymbolValue{Code: "USD", Value: 0.193584})
	})

	t.Run("Should return error when the API returns an error", func(t *testing.T) {
		srv := service.NewCurrencyService(baseURL, token, []string{"BRL"})

		result, err := srv.GetRate(context.Background(), "LLL")

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Cleanup(clearMock)
}
