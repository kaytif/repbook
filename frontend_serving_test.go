package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFrontendAndAPIRoutes(t *testing.T) {
	t.Chdir(t.TempDir())
	router := newRouter()
	for _, path := range []string{"/", "/api/workouts", "/workouts", "/api/missing"} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest("GET", path, nil))
		want := 401
		if path == "/" {
			want = 200
			if !strings.Contains(recorder.Body.String(), `<div id="root"></div>`) {
				t.Fatal("frontend index was not served")
			}
		} else if path == "/api/missing" {
			want = 404
		}
		if recorder.Code != want {
			t.Errorf("%s: got %d, want %d", path, recorder.Code, want)
		}
	}
}
