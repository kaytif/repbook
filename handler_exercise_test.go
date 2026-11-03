package main

import (
	"net/http"
	"testing"
)

func TestCreateExercise_Validation(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		body       string
		wantStatus int
	}{

		{"missing required fields (name)", createExercise, "POST", "/exercises", `{}`, 400},

		{"empty name", createExercise, "POST", "/exercises", `{"name":""}`, 400},

		{"whitespace name only", createExercise, "POST", "/exercises", `{"name":"  "}`, 400},

		{"malformed JSON", createExercise, "POST", "/exercises", `{"name":}`, 400},

		{"wrong data type", createExercise, "POST", "/exercises", `{"name":123}`, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := executeRequest(t, tt.handler, tt.method, tt.path, tt.body)

			if recorder.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d; response: %s",
					tt.wantStatus,
					recorder.Code,
					recorder.Body.String())
			}
		})
	}
}
