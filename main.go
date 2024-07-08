package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"starter/api"
	"starter/internal/database"
)

func main() {
	fmt.Println("Booting server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ResetDb()
	apiCfg := &api.ApiConfig{
		FileserverHits: 0,
		JwtSecret:      os.Getenv("JWT_SECRET"),
		PolkaApiKey:    os.Getenv("POLKA_API_KEY"),
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", HealthCheck)
	mux.HandleFunc("GET /admin/metrics", apiCfg.GetMetrics)
	mux.HandleFunc("/api/reset", apiCfg.ResetMetrics)

	mux.Handle("POST /api/chirps", apiCfg.RequireAuth(http.HandlerFunc(api.CreateChirp)))
	mux.HandleFunc("GET /api/chirps", api.GetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", api.GetChirp)
	mux.Handle("DELETE /api/chirps/{id}", apiCfg.RequireAuth(http.HandlerFunc(api.DeleteChirp)))

	mux.HandleFunc("POST /api/users", apiCfg.CreateUser)
	mux.Handle("PUT /api/users", apiCfg.RequireAuth(http.HandlerFunc(apiCfg.UpdateUser)))

	mux.HandleFunc("POST /api/login", apiCfg.Login)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeToken)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.WebhookHandler)

	handler := http.StripPrefix("/app", http.FileServer(http.Dir("")))
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(handler))

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}

func HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
