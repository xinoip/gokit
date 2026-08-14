package mux_test

import (
	"fmt"
	"net/http"

	"github.com/xinoip/gokit/mux"
)

func ExampleNewChi() {
	r := mux.NewChi(&mux.NewChiParams{
		AllowOrigins: []string{"https://app.example.com"},
	})
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	fmt.Println("health route configured")

	// Output:
	// health route configured
}
