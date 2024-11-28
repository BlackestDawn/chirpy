package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// prep
	godotenv.Load()

	// variables
	serverMux := http.NewServeMux()
	apiCfg, err := NewApiConfig()
	if err != nil {
		log.Fatalln(err)
	}
	handlerApp := http.StripPrefix("/app", http.FileServer(http.Dir(serverRootPath)))

	// path routing: general
	serverMux.Handle("/app/", apiCfg.middlewareMetricInc(handlerApp))

	// path routing: API
	serverMux.HandleFunc("GET /api/healthz", handlerHealth)

	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerNewChirp)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.handlerListChirps)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	serverMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	serverMux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	serverMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshAccessToken)
	serverMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)

	serverMux.HandleFunc("POST /api/users", apiCfg.handlerAddUser)
	serverMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	serverMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaUpgradeUser)

	// path routing: admin
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// initialization
	apiCfg.fileserverHits.Store(0)

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + serverListenPort,
	}

	// run it
	log.Println("Starting server")
	log.Fatal(server.ListenAndServe())
}
