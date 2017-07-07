package finalizer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleMiddleware() {
	//finalizer is executed at the end of the HTTP request in Middleware.
	finalizer := func(ctx context.Context, code int, r *http.Request) {
		header, _ := Header(ctx)
		size, _ := ResponseSize(ctx)
		contentType := header.Get("Content-Type")

		fmt.Printf("status=%d, contentType=%s, responseSize=%d\n", code, contentType, size)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Header().Set("Content-Type", "message/teapot")
		w.Write([]byte("I brew Tea"))
	})

	// wrap the handler with the finalizer middleware.
	wrappedHandler := Middleware(finalizer, handler)

	srv := httptest.NewServer(wrappedHandler)
	defer srv.Close()

	if _, err := http.Get(srv.URL); err != nil {
		panic(err)
	}

	// Output:
	// status=418, contentType=message/teapot, responseSize=0

}
