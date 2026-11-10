package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFrontendAndAPIRoutes(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.MkdirAll(filepath.Join("frontend", "dist"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("frontend", "dist", "index.html"), []byte("<html>Workout Tracker</html>"), 0644); err != nil {
		t.Fatal(err)
	}
	router := newRouter()
	for _, path := range []string{"/", "/api/workouts", "/workouts", "/api/missing"} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest("GET", path, nil))
		want := 401
		if path == "/" {
			want = 200
			if !strings.Contains(recorder.Body.String(), "Workout Tracker") {
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
