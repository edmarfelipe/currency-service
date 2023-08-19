package convert

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.come/edmarfelipe/currency-service/internal"
	"github.come/edmarfelipe/currency-service/internal/xhttp"
)

var (
	errCurrencyNotSupport   = xhttp.NewAPIError("Currency not supported", http.StatusBadRequest)
	errInvalidInputCurrency = xhttp.NewAPIError("Invalid 'Currency' parameter", http.StatusBadRequest)
	errInvalidInputValue    = xhttp.NewAPIError("Invalid 'Value' parameter", http.StatusBadRequest)
)

type controller struct {
	usc    *UseCase
	Config *internal.Config
}

func NewController(ct *internal.Container) xhttp.Controller {
	return &controller{
		Config: ct.Config,
		usc:    NewUseCase(ct.CurrencyService),
	}
}

func (c *controller) IsCurrencyIsSupport(currency string) bool {
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
		return errInvalidInputCurrency
	}

	value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 32)
	if err != nil || value == 0 {
		return errInvalidInputValue
	}

	if !c.IsCurrencyIsSupport(currency) {
		return errCurrencyNotSupport
	}

	in := Input{
		Currency: currency,
		Value:    value,
	}

	out, err := c.usc.Execute(r.Context(), in)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		return err
	}

	return nil
}
