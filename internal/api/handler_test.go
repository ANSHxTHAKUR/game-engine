package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ansh-singh/game-engine/internal/engine"
	"github.com/ansh-singh/game-engine/internal/model"
)

func TestSubmitHandler_GameActive(t *testing.T) {
	ge := engine.New(2, 64)
	defer ge.Shutdown()

	h := NewSubmitHandler(ge)

	body, _ := json.Marshal(model.UserResponse{
		UserID: "user_1", Answer: "42", IsCorrect: false,
	})
	req := httptest.NewRequest(http.MethodPost, "/submit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusAccepted)
	}

	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["status"] != "received" {
		t.Errorf("got status=%q, want %q", resp["status"], "received")
	}
}

func TestSubmitHandler_InvalidJSON(t *testing.T) {
	ge := engine.New(2, 64)
	defer ge.Shutdown()

	h := NewSubmitHandler(ge)

	req := httptest.NewRequest(http.MethodPost, "/submit", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestSubmitHandler_WrongMethod(t *testing.T) {
	ge := engine.New(2, 64)
	defer ge.Shutdown()

	h := NewSubmitHandler(ge)

	req := httptest.NewRequest(http.MethodGet, "/submit", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestMetricsHandler_ReturnsJSON(t *testing.T) {
	ge := engine.New(2, 64)
	defer ge.Shutdown()

	h := NewMetricsHandler(ge)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["correct"] != float64(0) {
		t.Errorf("got correct=%v, want 0", resp["correct"])
	}
	if resp["total"] != float64(0) {
		t.Errorf("got total=%v, want 0", resp["total"])
	}
	if _, ok := resp["uptime"]; !ok {
		t.Error("expected uptime field in response")
	}
}

func TestMetricsHandler_WrongMethod(t *testing.T) {
	ge := engine.New(2, 64)
	defer ge.Shutdown()

	h := NewMetricsHandler(ge)

	req := httptest.NewRequest(http.MethodPost, "/metrics", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
