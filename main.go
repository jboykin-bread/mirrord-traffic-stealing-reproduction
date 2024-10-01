package main

import (
	"flag"
	"log/slog"
	"net/http"
)

func main() {
	// flag for color
	color := flag.String("color", "blue", "the color to return in the response")
	flag.Parse()

	// http handler
	http.HandleFunc("/color", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"protocol", r.Proto,
			"remote_addr", r.RemoteAddr,
		)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(*color))
	})

	// start http server
	slog.Info("Starting server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		slog.Error("Error starting server", "error", err)
		panic(err)
	}
}
