package xhttp

import "net/http"

type Controller interface {
	Handler(w http.ResponseWriter, r *http.Request) error
}
