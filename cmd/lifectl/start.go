package main

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func startMultiverse(size int) (string, error) {
	pth, _ := url.JoinPath(host, "multiverse")
	pth += "?n=" + strconv.Itoa(size)
	resp, err := http.Post(pth, "", nil)
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
