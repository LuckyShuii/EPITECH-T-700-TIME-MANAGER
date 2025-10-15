package example

import "testing"

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
