package main

import (
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/kv"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		store, err := kv.OpenStore("default")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer store.Close()

		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)

		case http.MethodGet:
			w.WriteHeader(http.StatusOK)

		case http.MethodPut:
			w.WriteHeader(http.StatusOK)

		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)

		case http.MethodHead:
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func main() {}
