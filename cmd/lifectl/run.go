package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func runMultiverse() {
	unis, err := getUniverseList()
	if err != nil {
		fmt.Println(err)
		return
	}
	list := strings.Split(unis, "\n")
	var tc atomic.Int32

	go func() {
		for {
			var wg sync.WaitGroup
			for _, v := range list {
				if v == "" {
					continue
				}
				wg.Add(1)
				go func(id string, w *sync.WaitGroup) {
					defer w.Done()
					_, err := tickUniverse(id)
					if err != nil {
						fmt.Println(err)
						return
					}
					tc.Add(1)
				}(v, &wg)
			}
			wg.Wait()
		}
	}()

	// show snapshot of data/stats
	for {
		data, err := getUniverse(list[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
		ticks := tc.Swap(0)
		fmt.Println(float32(ticks)/float32(len(list)), "multiverses per second")
		fmt.Println(ticks, "universes per second")
		time.Sleep(1 * time.Second)
	}
}

func getUniverseList() (string, error) {
	pth, _ := url.JoinPath(host, "multiverse")
	resp, err := http.Get(pth)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func tickUniverse(id string) (string, error) {
	url, _ := url.JoinPath(host, "universe", id)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return "", errors.New("failed to tick universe")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getUniverse(id string) (string, error) {
	pth, _ := url.JoinPath(host, "universe", id)
	resp, err := http.Get(pth)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
