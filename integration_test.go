package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/ansh-singh/game-engine/internal/api"
	"github.com/ansh-singh/game-engine/internal/engine"
	"github.com/ansh-singh/game-engine/internal/model"
)

func TestIntegration_1000ConcurrentUsers(t *testing.T) {
	ge := engine.New(10, 2048)
	defer ge.Shutdown()

	ts := httptest.NewServer(api.NewSubmitHandler(ge))
	defer ts.Close()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			resp := model.UserResponse{
				UserID:    fmt.Sprintf("user_%d", id),
				Answer:    "42",
				IsCorrect: id%3 == 0,
			}
			body, _ := json.Marshal(resp)
			http.Post(ts.URL+"/submit", "application/json", bytes.NewReader(body))
		}(i)
	}

	select {
	case result := <-ge.Result():
		if result.WinnerID == "" {
			t.Error("expected a winner")
		}
		t.Logf("Winner: %s (correct=%d, incorrect=%d, time=%v)",
			result.WinnerID, result.TotalCorrect, result.TotalIncorrect, result.TimeTaken)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for winner")
	}

	wg.Wait()
}
