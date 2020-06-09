package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func HelloServer(w http.ResponseWriter, r *http.Request) {
	hlog.FromRequest(r).Info().
		Str("method", r.Method).
		Stringer("url", r.URL).
		Msg("")
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	http.HandleFunc("/", HelloServer)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
