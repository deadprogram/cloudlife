package main

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/acifani/vita/lib/game"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"github.com/fermyon/spin/sdk/go/v2/kv"
)

const (
	width, height  uint32 = 32, 32
	livePopulation        = 45
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
			key := game.GenerateKey()
			universe := game.NewDistributedUniverse(key, height, width)
			universe.Randomize(livePopulation)

			if err := store.Set(key, universe.Bytes()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(key)

		case http.MethodGet:
			key := path.Base(r.URL.Path)

			exists, err := store.Exists(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			value, err := store.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			universe := game.NewDistributedUniverse(key, height, width)
			universe.Write(value)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(universe.String()))

		case http.MethodPut:
			key := path.Base(r.URL.Path)

			exists, err := store.Exists(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			value, err := store.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			universe := game.NewDistributedUniverse(key, height, width)
			universe.Write(value)

			q := r.URL.Query()
			saveNeighbors := false
			if val := q.Get("topid"); val != "" {
				saveNeighbors = true
				universe.SetTopNeighbor(val)
			}
			if val := q.Get("bottomid"); val != "" {
				saveNeighbors = true
				universe.SetBottomNeighbor(val)
			}
			if val := q.Get("leftid"); val != "" {
				saveNeighbors = true
				universe.SetLeftNeighbor(val)
			}
			if val := q.Get("rightid"); val != "" {
				saveNeighbors = true
				universe.SetRightNeighbor(val)
			}

			if !saveNeighbors {
				universe.GetNeighbor = func(id string) *game.DistributedUniverse {
					value, err := store.Get(id)
					if err != nil {
						return nil
					}

					universe := game.NewDistributedUniverse(id, height, width)
					universe.Write(value)

					return universe
				}

				universe.Tick()
			}

			if err := store.Set(key, universe.Bytes()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			if saveNeighbors {
				json.NewEncoder(w).Encode(universe)
				return
			}
			w.Write([]byte(universe.String()))

		case http.MethodDelete:
			key := path.Base(r.URL.Path)

			exists, err := store.Exists(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !exists {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if err := store.Delete(key); err != nil {
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

func main() {}
