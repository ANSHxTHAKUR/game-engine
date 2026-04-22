package api

import (
	"encoding/json"
	"net/http"

	"github.com/ansh-singh/game-engine/internal/engine"
)

type MetricsHandler struct {
	engine engine.GameEngine
}

func NewMetricsHandler(ge engine.GameEngine) *MetricsHandler {
	return &MetricsHandler{engine: ge}
}

func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	m := h.engine.GetMetrics()
	correct, incorrect := m.Counts()
	elapsed := m.Elapsed()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"correct":   correct,
		"incorrect": incorrect,
		"total":     correct + incorrect,
		"uptime":    elapsed.String(),
	})
}
