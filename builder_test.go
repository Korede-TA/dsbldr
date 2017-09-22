package dsbldr

import (
	"reflect"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder(4, 100)
	if got, want := len(b.data), 101; got != want {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
	if got, want := len(b.data[0]), 4; got != want {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestAddDataFeature(t *testing.T) {
	b := NewBuilder(4, 3)
	b.addDataFeature("newFeature", []string{"one", "two", "three"})
	expectedData := [][]string{
		[]string{"newFeature", "", "", ""},
		[]string{"one", "", "", ""},
		[]string{"two", "", "", ""},
		[]string{"three", "", "", ""},
	}
	if got, want := b.data, expectedData; !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}
