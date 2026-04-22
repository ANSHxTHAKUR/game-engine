package model

import (
	"encoding/json"
	"testing"
)

func TestUserResponseJSON(t *testing.T) {
	raw := `{"user_id":"user_1","answer":"42","is_correct":true}`
	var resp UserResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.UserID != "user_1" {
		t.Errorf("got UserID=%q, want %q", resp.UserID, "user_1")
	}
	if resp.Answer != "42" {
		t.Errorf("got Answer=%q, want %q", resp.Answer, "42")
	}
	if !resp.IsCorrect {
		t.Error("got IsCorrect=false, want true")
	}
}

func TestUserResponseJSONIncorrect(t *testing.T) {
	raw := `{"user_id":"user_2","answer":"wrong","is_correct":false}`
	var resp UserResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.IsCorrect {
		t.Error("got IsCorrect=true, want false")
	}
}
