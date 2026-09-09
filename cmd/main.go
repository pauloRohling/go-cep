package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	BrasilAPI = "https://brasilapi.com.br/api/cep/v1/%s"
	ViaCepAPI = "http://viacep.com.br/ws/%s/json"
)

type Result struct {
	response string
	api      string
}

func main() {
	if len(os.Args) < 2 {
		panic("Usage: go run cmd/main.go <CEP>")
	}

	inputCEP := os.Args[1]

	brasilApiUrl := fmt.Sprintf(BrasilAPI, inputCEP)
	viaCepUrl := fmt.Sprintf(ViaCepAPI, inputCEP)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resultChan := make(chan Result)

	go func() {
		result, _ := Get(ctx, brasilApiUrl)
		resultChan <- Result{response: result, api: "BrasilAPI"}
	}()

	go func() {
		result, _ := Get(ctx, viaCepUrl)
		resultChan <- Result{response: result, api: "ViaCepAPI"}
	}()

	select {
	case r := <-resultChan:
		cancel()
		fmt.Printf("%s: %s", r.api, r.response)
		return
	case <-ctx.Done():
		fmt.Printf("Timeout!\n")
		return
	}
}

func Get(ctx context.Context, url string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
