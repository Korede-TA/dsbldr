package dsbldr

import (
	"testing"
)

func TestGetParentNames(t *testing.T) {
	f := NewFeature()
	f.Endpoint = "/items/{{item_id}}/categories"

	got, err := f.getParentNames()
	if err != nil {
		t.Errorf("Error Occured: %v", err)
	}

	if want := []string{"item_id"}; got[0] != want[0] {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}
