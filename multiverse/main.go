package main

import (
	"errors"
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

			for i := 0; i < defaultCount*defaultCount; i++ {
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
				universes = append(universes, strings.Trim(string(body), "\n\""))
			}

			if err := connectUniverses(universes); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(universes, "\n")))
			w.Write([]byte("\n"))

		case http.MethodGet:
			// get list of the universes
			unis, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(unis, "\n")))
			w.Write([]byte("\n"))

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

func connectUniverses(multi []string) error {
	number := defaultCount
	i := 0
	for row := 0; row < number; row++ {
		for col := 0; col < number; col++ {
			switch {
			// top right
			case col == number-1 && row == number-1:
				err := updateNeighbors(multi[i], "", multi[i-number], multi[i-1], "")
				if err != nil {
					return err
				}

			// top left
			case col == 0 && row == number-1:
				err := updateNeighbors(multi[i], "", multi[i-number], "", multi[i+1])
				if err != nil {
					return err
				}

			// bottom right
			case col == number-1 && row == 0:
				err := updateNeighbors(multi[i], multi[i+number], "", multi[i-1], "")
				if err != nil {
					return err
				}

			// bottom left
			case col == 0 && row == 0:
				err := updateNeighbors(multi[i], multi[i+number], "", "", multi[i+1])
				if err != nil {
					return err
				}

			// first column
			case col == 0:
				err := updateNeighbors(multi[i], multi[i+number], multi[i-number], "", multi[i+1])
				if err != nil {
					return err
				}

			// first row
			case row == 0:
				err := updateNeighbors(multi[i], multi[i+number], "", multi[i-1], multi[i+1])
				if err != nil {
					return err
				}

			// last column
			case col == number-1:
				err := updateNeighbors(multi[i], multi[i+number], multi[i-number], multi[i-1], "")
				if err != nil {
					return err
				}

			// last row
			case row == number-1:
				err := updateNeighbors(multi[i], "", multi[i-number], multi[i-1], multi[i+1])
				if err != nil {
					return err
				}

			// anyplace else
			default:
				err := updateNeighbors(multi[i], multi[i+number], multi[i-number], multi[i-1], multi[i+1])
				if err != nil {
					return err
				}
			}
			i++
		}
	}
	return nil
}

func updateNeighbors(id, top, bottom, left, right string) error {
	pth, _ := url.JoinPath("/universe", id)
	u := &url.URL{
		Path: pth,
	}
	u.RawQuery = "topid=" + top + "&bottomid=" + bottom + "&leftid=" + left + "&rightid=" + right

	req := http.Request{
		Method: http.MethodPut,
		URL:    u,
	}
	r, err := spinhttp.Send(&req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return errors.New("failed to update neighbors")
	}

	return nil
}

func main() {}
