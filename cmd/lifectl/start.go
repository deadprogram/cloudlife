package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func startMultiverse(size int) (string, error) {
	pth, _ := url.JoinPath(host, "multiverse")
	pth += "?n=" + strconv.Itoa(size)
	r, err := http.Post(pth, "", nil)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return "", errors.New("failed to start multiverse")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
