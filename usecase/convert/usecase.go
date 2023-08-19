package convert

import (
	"context"
	"fmt"
	"log/slog"

	"github.come/edmarfelipe/currency-service/internal/extsrv"
)

type Input struct {
	Currency string
	Value    float64
}

type CurrencyValue struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

func NewCurrencyValue(currency string, value float64) CurrencyValue {
	return CurrencyValue{
		Currency: currency,
		Value:    fmt.Sprintf("%.4f", value),
	}
}

type UseCase struct {
	service extsrv.CurrencyService
}

func NewUseCase(service extsrv.CurrencyService) *UseCase {
	return &UseCase{
		service: service,
	}
}

func (u *UseCase) Execute(ctx context.Context, input Input) ([]CurrencyValue, error) {
	slog.InfoContext(ctx, "Executing use case convert", "input", input)
	rate, err := u.service.GetRate(ctx, input.Currency)
	if err != nil {
		return nil, err
	}

	var output []CurrencyValue
	for _, r := range rate {
		if r.Code == input.Currency {
			continue
		}
		output = append(output, NewCurrencyValue(r.Code, input.Value*r.Value))
	}
	return output, nil
}
