package extsrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var (
	ErrFailedToConnect = errors.New("failed to connect to external service")
)

// CurrencyService is an interface for the external currency service.
type CurrencyService interface {
	// GetRate gets the rate of the given currency to the currencies defined in the service.
	GetRate(ctx context.Context, currency string) ([]SymbolValue, error)
}

// SymbolValue is a struct that represents the rate of a currency.
type SymbolValue struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

type apiResult struct {
	Meta struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `json:"meta"`
	Data map[string]SymbolValue `json:"data"`
}

type currencyService struct {
	apiURL     string
	apiToken   string
	currencies []string
	httpClient *http.Client
}

// NewCurrencyService creates a new currency service.
func NewCurrencyService(apiURL string, apiToken string, currencies []string) CurrencyService {
	return &currencyService{
		apiURL:     apiURL,
		apiToken:   apiToken,
		currencies: currencies,
		httpClient: &http.Client{
			Timeout: time.Second * 2,
		},
	}
}

func (srv *currencyService) GetRate(ctx context.Context, currency string) ([]SymbolValue, error) {
	url := fmt.Sprintf("%s?base_currency=%s&currencies=%s", srv.apiURL, currency, strings.Join(srv.currencies, ","))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("apikey", srv.apiToken)
	resp, err := srv.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			slog.ErrorContext(ctx, "Could not close body response", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrFailedToConnect
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response apiResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var result []SymbolValue
	for _, value := range response.Data {
		result = append(result, value)
	}

	return result, nil
}
