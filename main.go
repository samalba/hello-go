package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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
	if r.URL.RawQuery == "debug" {
		fmt.Fprintf(w, "\n\n\n-- env:\n\n")
		for _, e := range os.Environ() {
			fmt.Fprintf(w, "%s\n", e)
		}
		fmt.Fprintf(w, "\n\n-- headers:\n\n")
		for k, v := range r.Header {
			fmt.Fprintf(w, "%s: %s\n", k, strings.Join(v, ","))
		}
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal().Msg("PORT is undefined")
	}

	http.HandleFunc("/", HelloServer)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
