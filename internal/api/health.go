package api

import "net/http"

func (s *Server) liveness(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) readiness(w http.ResponseWriter, r *http.Request) {
	if err := s.db.PingContext(r.Context()); err != nil {
		s.log.Error("readiness probe failed", "error", err)
		WriteError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
