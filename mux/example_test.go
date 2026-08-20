package mux_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/xinoip/gokit/mux"
)

func ExampleNewChi() {
	r := mux.NewChi(mux.DefaultChiConfig())
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		_, err := io.WriteString(w, "healthy")
		if err != nil {
			panic(err)
		}
	})

	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
	request.Header.Set("Origin", "https://app.example.com")
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	fmt.Println(response.Code)
	fmt.Println(response.Header().Get("Access-Control-Allow-Origin"))
	fmt.Println(response.Body.String())

	// Output:
	// 200
	// https://app.example.com
	// healthy
}
