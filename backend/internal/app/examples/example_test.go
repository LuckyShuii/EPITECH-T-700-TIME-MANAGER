package example

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ExampleHandler "app/internal/app/examples/handler"

	"github.com/gin-gonic/gin"
)

func TestGreet(t *testing.T) {
	got := Greet("Lucas")
	want := "Hello, Lucas!"
	if got != want {
		t.Fatalf("Greet(\"Lucas\") = %q; want %q", got, want)
	}

	got = Greet("")
	want = "Hello, world!"
	if got != want {
		t.Fatalf("Greet(\"\") = %q; want %q", got, want)
	}
}

func TestHelloHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/hello", ExampleHandler.ExampleHandler)

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d; want %d", w.Code, http.StatusOK)
	}

	expected := `{"message":"hello"}`
	if w.Body.String() != expected {
		t.Fatalf("body = %q; want %q", w.Body.String(), expected)
	}
}
