package main

import (
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
			key := generateKey()
			universe := game.NewUniverse(height, width)
			universe.Randomize(livePopulation)

			data := NewDataRecord(key)
			if _, err := universe.Read(data.Cells); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			buf := StoreFromDataRecord(data)

			if err := store.Set(key, buf); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(key))

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

			data := DataRecordFromStore(value)
			universe := UniverseFromDataRecord(data)

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

			data := DataRecordFromStore(value)
			universe := UniverseFromDataRecord(data)

			universe.Tick()

			universe.Read(data.Cells)

			buf := StoreFromDataRecord(data)

			if err := store.Set(key, buf); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
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
