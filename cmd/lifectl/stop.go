package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

func stopMultiverse() (string, error) {
	url, _ := url.JoinPath(host, "multiverse")
	req, err := http.NewRequest(http.MethodDelete, url, nil)
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
		return "", errors.New("failed to stop multiverse")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
