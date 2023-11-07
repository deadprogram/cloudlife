package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

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
		switch r.Method {
		case http.MethodPost:
			store, err := kv.OpenStore("default")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer store.Close()

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
