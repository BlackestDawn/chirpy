package main

import (
	"log"
	"net/http"
)

func main() {
	// variables
	serverMux := http.NewServeMux()
	apiCfg := NewApiConfig()
	handlerApp := http.StripPrefix("/app", http.FileServer(http.Dir(serverRootPath)))

	// path routing: general
	serverMux.Handle("/app/", apiCfg.middlewareMetricInc(handlerApp))

	// path routing: API
	serverMux.HandleFunc("GET /api/healthz", handlerHealth)

	// path routing: admin
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.resetHits)

	// initialization
	apiCfg.fileserverHits.Store(0)

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + serverListenPort,
	}

	// run it
	log.Fatal(server.ListenAndServe())
}
