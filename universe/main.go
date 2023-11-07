package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"path"

	"github.com/acifani/vita/lib/game"
	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/kv"
)

const (
	width, height  uint32 = 32, 32
	livePopulation        = 75
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
			key := generateKey()
			universe := game.NewUniverse(height, width)
			universe.Randomize(livePopulation)

			data := make([]byte, height*width)
			universe.Read(data)

			if err := store.Set(key, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(key))

		case http.MethodGet:
			value, err := store.Get(path.Base(r.URL.Path))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			universe := game.NewUniverse(height, width)
			universe.Write(value)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(universe.String()))

		case http.MethodPut:
			key := path.Base(r.URL.Path)
			value, err := store.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			universe := game.NewUniverse(height, width)
			universe.Write(value)
			universe.Tick()

			data := make([]byte, height*width)
			universe.Read(data)

			if err := store.Set(key, data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(universe.String()))

		case http.MethodDelete:
			if err := store.Delete(path.Base(r.URL.Path)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)

		case http.MethodHead:
			exists, err := store.Exists(path.Base(r.URL.Path))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if exists {
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusNotFound)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func generateKey() string {
	var result [32]byte
	rand.Read(result[:])
	encodedString := hex.EncodeToString(result[:])
	return encodedString
}

func main() {}
