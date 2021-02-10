package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/chyeh/pubip"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	if r.URL.RawQuery == "ip" {
		fmt.Fprintf(w, "\n\n\n-- my public ip: ")
		ip, err := pubip.Get()
		if err != nil {
			fmt.Fprintf(w, "%s\n", err)
		} else {
			fmt.Fprintf(w, "%s\n", ip)
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Warn().Msg("PORT is undefined, assuming 80")
		port = "80"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/*", HelloServer)
	r.Get("/db/{dbPrefix}", TestDB)
	fmt.Println("Listening to port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
