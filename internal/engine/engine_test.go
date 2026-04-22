package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/ansh-singh/game-engine/internal/model"
)

func TestNewGameEngine(t *testing.T) {
	ge := New(10, 1024)
	defer ge.Shutdown()
}

func TestSubmitAfterShutdown(t *testing.T) {
	ge := New(10, 1024)
	ge.Shutdown()

	err := ge.Submit(model.UserResponse{
		UserID:    "user_1",
		Answer:    "42",
		IsCorrect: true,
		Timestamp: time.Now(),
	})
	if err == nil {
		t.Error("expected error on submit after shutdown")
	}
}

func TestSingleCorrectAnswer(t *testing.T) {
	ge := New(4, 64)
	defer ge.Shutdown()

	ge.Submit(model.UserResponse{
		UserID: "user_1", Answer: "42", IsCorrect: true, Timestamp: time.Now(),
	})

	select {
	case result := <-ge.Result():
		if result.WinnerID != "user_1" {
			t.Errorf("got winner=%q, want %q", result.WinnerID, "user_1")
		}
		if result.TotalCorrect != 1 {
			t.Errorf("got correct=%d, want 1", result.TotalCorrect)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for result")
	}
}

func TestOnlyOneWinner(t *testing.T) {
	ge := New(4, 128)
	defer ge.Shutdown()

	for i := 0; i < 100; i++ {
		ge.Submit(model.UserResponse{
			UserID:    fmt.Sprintf("user_%d", i),
			Answer:    "42",
			IsCorrect: true,
			Timestamp: time.Now(),
		})
	}

	select {
	case result := <-ge.Result():
		if result.WinnerID == "" {
			t.Error("expected a winner, got empty string")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for result")
	}

	select {
	case extra := <-ge.Result():
		t.Errorf("got second result: %+v — only one winner should be declared", extra)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestIncorrectAnswersNoWinner(t *testing.T) {
	ge := New(4, 64)
	defer ge.Shutdown()

	for i := 0; i < 50; i++ {
		ge.Submit(model.UserResponse{
			UserID: fmt.Sprintf("user_%d", i), Answer: "wrong", IsCorrect: false, Timestamp: time.Now(),
		})
	}

	select {
	case result := <-ge.Result():
		t.Errorf("should not have a winner with all incorrect answers, got: %+v", result)
	case <-time.After(200 * time.Millisecond):
	}
}
