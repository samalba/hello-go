package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/hlog"
)

func readDBCreds(w http.ResponseWriter, dbPrefix string) (string, string) {
	keyPrefix := strings.ToUpper(dbPrefix)
	keysList := []string{
		"DBNAME",
		"HOSTNAME",
		"USERNAME",
		"PASSWORD",
		"PORT",
		"TYPE",
	}
	keys := map[string]string{}
	for _, k := range keysList {
		envKey := keyPrefix + k
		envVar := os.Getenv(envKey)
		if envVar == "" {
			fmt.Fprintf(w, "Warning: %q env var is empty\n", envKey)
		}
		keys[k] = envVar
	}
	return keys["TYPE"], fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		keys["TYPE"],
		keys["USERNAME"],
		keys["PASSWORD"],
		keys["HOSTNAME"],
		keys["PORT"],
		keys["DBNAME"])
}

func TestDB(w http.ResponseWriter, r *http.Request) {
	hlog.FromRequest(r).Info().
		Str("method", r.Method).
		Stringer("url", r.URL).
		Msg("")
	dbType, connStr := readDBCreds(w, chi.URLParam(r, "dbPrefix"))
	if dbType != "mysql" && dbType != "postgres" {
		fmt.Fprintf(w, "Error: unknown db type %q\n", dbType)
		return
	}
	fmt.Fprintf(w, "sql.Open(%q, %q) => ", dbType, connStr)
	db, err := sql.Open(dbType, connStr)
	if err != nil {
		fmt.Fprintf(w, "ERROR: %s\n", err)
	} else {
		fmt.Fprintf(w, "OK\n")
	}
	fmt.Fprintf(w, "db.Ping() => ")
	if err := db.Ping(); err != nil {
		fmt.Fprintf(w, "ERROR: %s\n", err)
	} else {
		fmt.Fprintf(w, "OK\n")
	}
	stats := db.Stats()
	statsJson, err := json.MarshalIndent(stats, "", "    ")
	if err != nil {
		fmt.Fprintf(w, "db.Stats() => %s\n", statsJson)
	}
	fmt.Fprintf(w, "db.Close() => ")
	if err := db.Close(); err != nil {
		fmt.Fprintf(w, "ERROR: %s\n", err)
	} else {
		fmt.Fprintf(w, "OK\n")
	}
}
