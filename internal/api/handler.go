package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ansh-singh/game-engine/internal/engine"
	"github.com/ansh-singh/game-engine/internal/model"
)

type SubmitHandler struct {
	engine engine.GameEngine
}

func NewSubmitHandler(ge engine.GameEngine) *SubmitHandler {
	return &SubmitHandler{engine: ge}
}

func (h *SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var resp model.UserResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	resp.Timestamp = time.Now()

	if err := h.engine.Submit(resp); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "game_over",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "received",
	})
}
