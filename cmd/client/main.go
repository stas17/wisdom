package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"wisdom/internal/pkg/pow"
)

const method = "http://wisdom:8080/wisdom"

func main() {
	resp, err := http.Get(method)
	if err != nil {
		log.Fatalf("Send error: %s", err)
	}
	var res string

	json.NewDecoder(resp.Body).Decode(&res)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Fail. Code: %d, Msg: %s", resp.StatusCode, res)
	}

	block := pow.NewBlock(res)

	body, err := json.Marshal(block)
	if err != nil {
		log.Fatalf("Marshal error: %s", err)
	}

	reqURL, err := url.Parse(method)
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"Salt":         {res},
		},
		Body: ioutil.NopCloser(bytes.NewReader(body)),
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Send error: %s", err)
	}
	json.NewDecoder(resp.Body).Decode(&res)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Fail. Code: %d, Msg: %s", resp.StatusCode, res)
	}

	fmt.Println(res)
}
