package main

import (
	"net/http"
	"testing"
)

func TestCreateProgram_Validation(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		method     string
		path       string
		body       string
		wantStatus int
	}{

		{"missing required fields (name)", createProgram, "POST", "/programs", `{}`, 400},

		{"empty name", createProgram, "POST", "/programs", `{"name":""}`, 400},

		{"whitespace name only", createProgram, "POST", "/programs", `{"name":"  "}`, 400},

		{"malformed JSON", createProgram, "POST", "/programs", `{"name":}`, 400},

		{"wrong data type", createProgram, "POST", "/programs", `{"name":123}`, 400},
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
