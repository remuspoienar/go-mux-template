package api

import (
	"fmt"
	"net/http"
)

func (c *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.FileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (c *ApiConfig) GetMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>`, c.FileserverHits)
	w.Write([]byte(html))
}

func (c *ApiConfig) ResetMetrics(w http.ResponseWriter, req *http.Request) {
	c.FileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", c.FileserverHits)))
}
