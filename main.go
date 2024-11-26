package main

import (
	"log"
	"net/http"
)

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	serverMux := http.NewServeMux()
	apiCfg := NewApiConfig()
	handlerApp := http.StripPrefix("/app", http.FileServer(http.Dir(serverRootPath)))
	serverMux.Handle("/app/", apiCfg.middlewareMetricInc(handlerApp))
	serverMux.HandleFunc("GET /healthz", handlerHealth)
	serverMux.HandleFunc("GET /metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("POST /reset", apiCfg.resetHits)
	apiCfg.fileserverHits.Store(0)

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + serverListenPort,
	}

	log.Fatal(server.ListenAndServe())
}
