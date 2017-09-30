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

func TestResolveEndpoint(t *testing.T) {
	f := NewFeature()
	f.Endpoint = "/items/category/{{item_ids}}/{{other}}"

	parentsValues := map[string]string{
		"item_ids": "item1123",
		"other":    "anotherthing",
	}

	got, err := f.resolveEndpoint(parentsValues)

	if err != nil {
		t.Errorf("Error Occured: %v", err)
	}

	if want := "/items/category/item1123/anotherthing"; got != want {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}
