package api

import "net/http"

// @Summary      Liveness probe
// @Description  Returns 200 if the process is running.
// @Tags         health
// @Produce      json
// @Success      200  {object}  healthResponse
// @Router       /healthz/live [get]
func (s *Server) liveness(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// @Summary      Readiness probe
// @Description  Returns 200 if the database is reachable, 503 otherwise.
// @Tags         health
// @Produce      json
// @Success      200  {object}  healthResponse
// @Failure      503  {object}  errorResponse
// @Router       /healthz/ready [get]
func (s *Server) readiness(w http.ResponseWriter, r *http.Request) {
	if err := s.db.PingContext(r.Context()); err != nil {
		s.log.Error("readiness probe failed", "error", err)
		WriteError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
