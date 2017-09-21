package dsbldr

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder(4, 100)
	if got, want := len(b.data), 100; got != want {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
	if got, want := len(b.data[0]), 4; got != want {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}
