package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/kv"
)

const defaultCount = 4

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
			// create all of the universes
			var universes []string

			for i := 0; i < defaultCount; i++ {
				r, err := spinhttp.Post("/universe", "application/text", nil)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer r.Body.Close()

				if r.StatusCode != http.StatusOK {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				body, err := io.ReadAll(r.Body)
				universes = append(universes, string(body))
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(universes, "/n")))

		case http.MethodGet:
			// get list of the universes
			unis, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(unis, "/n")))

		case http.MethodPut:
			// run a generation on all universes
			unis, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var universes []string
			for _, v := range unis {
				pth, _ := url.JoinPath("/universe", v)
				u := &url.URL{
					Path: pth,
				}

				req := http.Request{
					Method: http.MethodPut,
					URL:    u,
				}
				r, err := spinhttp.Send(&req)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer r.Body.Close()

				if r.StatusCode != http.StatusOK {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				body, err := io.ReadAll(r.Body)
				universes = append(universes, string(body))
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(universes, "")))

		case http.MethodDelete:
			// delete all universes
			unis, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, v := range unis {
				pth, _ := url.JoinPath("/universe", v)
				u := &url.URL{
					Path: pth,
				}

				req := http.Request{
					Method: http.MethodDelete,
					URL:    u,
				}
				r, err := spinhttp.Send(&req)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if r.StatusCode != http.StatusOK {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusOK)

		case http.MethodHead:
			// TODO: check if any universes out there?

			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func main() {}
