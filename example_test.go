package slogio_test

import (
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/utkuozdemir/go-slogio"
)

func ExampleWriter() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Time(slog.TimeKey, time.Unix(0, 0).UTC())
			}

			return a
		},
	}))

	w := &slogio.Writer{Log: logger}

	io.WriteString(w, "starting up\n")
	io.WriteString(w, "running\n")
	io.WriteString(w, "shutting down\n")

	if err := w.Close(); err != nil {
		log.Fatal(err)
	}

	// Output:
	// {"time":"1970-01-01T00:00:00Z","level":"INFO","msg":"starting up"}
	// {"time":"1970-01-01T00:00:00Z","level":"INFO","msg":"running"}
	// {"time":"1970-01-01T00:00:00Z","level":"INFO","msg":"shutting down"}
}
