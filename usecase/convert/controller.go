package convert

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.come/edmarfelipe/currency-service/internal"
)

var (
	ErrCurrencyNotSupport   = internal.APIError{Message: "Currency not supported", Status: http.StatusBadRequest}
	ErrInvalidInputCurrency = internal.APIError{Message: "Invalid 'Currency' parameter", Status: http.StatusBadRequest}
	ErrInvalidInputValue    = internal.APIError{Message: "Invalid 'Value' parameter", Status: http.StatusBadRequest}
)

type controller struct {
	usc    *UseCase
	Config *internal.Config
}

func NewController(ct *internal.Container) *controller {
	return &controller{
		Config: ct.Config,
		usc:    NewUseCase(ct.CurrencyService),
	}
}

func (c *controller) isCurrencyIsSupport(currency string) bool {
	for _, item := range c.Config.Currencies {
		if item == currency {
			return true
		}
	}
	return false
}

func (c *controller) Handler(w http.ResponseWriter, r *http.Request) error {
	currency := chi.URLParam(r, "currency")
	if len(currency) < 3 || len(currency) > 3 {
		return ErrInvalidInputCurrency
	}

	value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 32)
	if err != nil || value == 0 {
		return ErrInvalidInputValue
	}

	if !c.isCurrencyIsSupport(currency) {
		return ErrCurrencyNotSupport
	}

	in := Input{
		Currency: currency,
		Value:    value,
	}

	out, err := c.usc.Execute(r.Context(), in)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(out)
}
