package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func do(verb, host, path, contentType string, input []byte,
	decode func(http.Header, io.Reader) error) error {

	u := fmt.Sprintf("http://%s%s", host, path)
	rq, err := http.NewRequest(verb, u, bytes.NewBuffer(input))
	if err != nil {
		return err
	}
	rq.Header.Set("Content-Type", contentType)

	client := http.Client{}
	rsp, err := client.Do(rq)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	var body io.Reader = rsp.Body

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: %v status code: %v", u, rsp.StatusCode)
	}
	if decode == nil {
		return nil
	}
	return decode(rsp.Header, body)
}

func post(host, path string, input interface{}, reader func(http.Header, io.Reader) error) error {
	buf, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return do("POST", host, path, "application/json", buf, reader)
}

func put(host, path string, input interface{}, reader func(http.Header, io.Reader) error) error {
	buf, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return do("PUT", host, path, "application/json", buf, reader)
}

func get(host, path string, read func(http.Header, io.Reader) error) error {
	return do("GET", host, path, "application/json", nil, read)
}

func postJson(host, path string, input, output interface{}) error {
	return post(host, path, input,
		func(_ http.Header, r io.Reader) error {
			return json.NewDecoder(r).Decode(output)
		})
}

func putJson(host, path string, input, output interface{}) error {
	return put(host, path, input,
		func(_ http.Header, r io.Reader) error {
			return json.NewDecoder(r).Decode(output)
		})
}

func getJson(host, path string, output interface{}) error {
	return get(host, path,
		func(_ http.Header, r io.Reader) error {
			return json.NewDecoder(r).Decode(output)
		})
}
