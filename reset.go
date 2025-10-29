package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hits reset to 0")
}