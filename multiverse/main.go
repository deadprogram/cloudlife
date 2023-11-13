package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
			existing, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if len(existing) > 0 {
				http.Error(w, "universes already exist in multiverse", http.StatusConflict)
				return
			}

			n := defaultCount
			q := r.URL.Query()
			if val := q.Get("n"); val != "" {
				if c, err := strconv.Atoi(val); err == nil {
					n = c
				}
			}

			universes, err := createUniverses(n)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := connectUniverses(universes, n); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(universes, "\n")))
			w.Write([]byte("\n"))

		case http.MethodGet:
			// get keys to the universes
			unis, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(unis, "\n")))
			w.Write([]byte("\n"))

		case http.MethodPut:
			// run tick on all universes
			unis, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			result, err := multitick(unis)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(result, "")))

		case http.MethodDelete:
			// delete all universes
			unis, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := deleteUniverses(unis); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)

		case http.MethodHead:
			_, err := store.GetKeys()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func createUniverses(n int) ([]string, error) {
	var unis []string
	for i := 0; i < n*n; i++ {
		u, err := createUniverse()
		if err != nil {
			return unis, err
		}

		unis = append(unis, u)
	}

	return unis, nil
}

func createUniverse() (string, error) {
	r, err := spinhttp.Post("/universe", "application/text", nil)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return "", errors.New("failed to create universe")
	}

	body, err := io.ReadAll(r.Body)
	return strings.Trim(string(body), "\n\""), nil
}

func connectUniverses(multi []string, n int) error {
	i := 0
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			var err error

			switch {
			// top right
			case col == n-1 && row == n-1:
				err = updateNeighbors(multi[i], "", multi[i-n], multi[i-1], "")

			// top left
			case col == 0 && row == n-1:
				err = updateNeighbors(multi[i], "", multi[i-n], "", multi[i+1])

			// bottom right
			case col == n-1 && row == 0:
				err = updateNeighbors(multi[i], multi[i+n], "", multi[i-1], "")

			// bottom left
			case col == 0 && row == 0:
				err = updateNeighbors(multi[i], multi[i+n], "", "", multi[i+1])

			// first column
			case col == 0:
				err = updateNeighbors(multi[i], multi[i+n], multi[i-n], "", multi[i+1])

			// first row
			case row == 0:
				err = updateNeighbors(multi[i], multi[i+n], "", multi[i-1], multi[i+1])

			// last column
			case col == n-1:
				err = updateNeighbors(multi[i], multi[i+n], multi[i-n], multi[i-1], "")

			// last row
			case row == n-1:
				err = updateNeighbors(multi[i], "", multi[i-n], multi[i-1], multi[i+1])

			// anyplace else
			default:
				err = updateNeighbors(multi[i], multi[i+n], multi[i-n], multi[i-1], multi[i+1])
			}

			if err != nil {
				return err
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

func multitick(unis []string) ([]string, error) {
	result := []string{}
	for _, v := range unis {
		u, err := tick(v)
		if err != nil {
			return result, err
		}

		result = append(result, u)
	}

	return result, nil
}

func tick(key string) (string, error) {
	pth, _ := url.JoinPath("/universe", key)
	u := &url.URL{
		Path: pth,
	}

	req := http.Request{
		Method: http.MethodPut,
		URL:    u,
	}
	r, err := spinhttp.Send(&req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return "", errors.New("failed to advance universes")
	}

	body, err := io.ReadAll(r.Body)
	return string(body), nil
}

func deleteUniverses(unis []string) error {
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
			return err
		}

		if r.StatusCode != http.StatusOK {
			return errors.New("failed to delete universes")
		}
	}

	return nil
}

func main() {}
