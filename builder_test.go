package dsbldr

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
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
	if got, want := b.records, 100; got != want {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestAddFeatureData(t *testing.T) {
	b := NewBuilder(4, 3)
	b.addFeatureData("newFeature", []string{"one", "two", "three"})
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

func TestAddFeature(t *testing.T) {
	f := &Feature{
		Name:     "feat1",
		Endpoint: "/endpoint1/",
		RunFunc: func(res []string) []string {
			return []string{"one", "two", "three"}
		},
	}
	b := NewBuilder(4, 3)
	b.AddFeatures(f)
	if got, want := b.featureMap["feat1"], f; got != want {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestGetFeatureData(t *testing.T) {
	f := &Feature{
		Name:     "feat1",
		Endpoint: "/endpoint1/",
	}
	// Note that the test fails when there is a greater featureCount
	// than there are features when builder.getFeatureData is called
	b := NewBuilder(1, 3)
	data := []string{"one", "two", "three"}

	t.Log(b.data)
	err := b.addFeatureData(f.Name, data)
	t.Log(b.data)
	if err != nil {
		t.Errorf("Error Occured: %v", err)
	}

	want := []string{"one", "two", "three"}
	if got := b.getFeatureData("feat1"); !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestGetFeature(t *testing.T) {
	f := &Feature{
		Name:     "feat1",
		Endpoint: "/endpoint1/",
		RunFunc: func(res []string) []string {
			return []string{"one", "two", "three"}
		},
	}
	b := NewBuilder(4, 3)
	b.AddFeatures(f)
	if got, want := b.GetFeature("feat1"), f; !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestGetFeatureIfNoFeatures(t *testing.T) {
	var NilFeature *Feature
	b := NewBuilder(4, 3)
	if got, want := b.GetFeature("feat1"), NilFeature; !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestCreateRequest(t *testing.T) {
	b := NewBuilder(4, 3) // arbitrary preallocations
	b.BaseURL = "baseurl.com"
	sampleOAuthHeader := BasicOAuthHeader(
		"consumerKey", "nonce", "signature", "signatureMethod",
		"timestamp", "token",
	)
	b.RequestHeaders = map[string]string{
		"OAuth": sampleOAuthHeader,
	}
	endpoint := "/myEndpoint"

	got, err := b.createRequest(endpoint)
	if err != nil {
		t.Errorf("Error Occured: %v", err)
	}

	u, err := url.Parse(b.BaseURL + endpoint)
	if err != nil {
		t.Errorf("Error Occured: %v", err)
	}

	want := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: map[string][]string{
			"Authorization": []string{sampleOAuthHeader},
		},
		Body: nil,
		Host: u.Host,
	}

	if reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}

}

func TestResolveFeatureEndpoints(t *testing.T) {
	b := NewBuilder(2, 3)
	b.BaseURL = "baseurl.com"

	f1 := &Feature{
		Name:     "feat1",
		Endpoint: "/endpoint1/",
	}

	f2 := &Feature{
		Name:     "feat2",
		Endpoint: "/endpoint2/{{feat1}}",
	}

	b.AddFeatures(f1, f2)
	b.addFeatureData(f1.Name, []string{"one", "two", "three"})

	got, err := b.resolveFeatureEndpoints(f2)
	if err != nil {
		t.Errorf("Error Occured: %v", err)
	}

	want := []string{
		"/endpoint2/one",
		"/endpoint2/two",
		"/endpoint2/three",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

// func TestRetrieveFeature(t *testing.T) {
// 	b := NewBuilder(4, 3)
// }

type fakeHttpClient struct{}

var fakeResponseDump string = `
{
	id: 1000
	id_str: "1000"
	text: "this is a young tweet; tweety tweet tweet"
}
`

func (fhc *fakeHttpClient) Do(req http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(fakeResponseDump)),
	}, nil
}

func TestPopulateFeatureData(t *testing.T) {
	b := NewBuilder(2, 3)
	fakeClient := fakeHttpClient{}
	b.BaseURL = "baseurl.com"
	f := &Feature{
		Name:     "f1",
		Endpoint: "/endpoint",
		RunFunc: func(res []string) []string {
			return []string{"one", "two", "three"}
		},
	}

	got, err := b.populateFeatureData(f, &fakeClient)
	if err != nil {
		t.Errorf("Error Occured: %v", err)
	}
	if want := fakeResponseDump; got[0] != want {
		t.Fatalf("got: %v\n want: %v\n", got, want)
	}
}

func TestRun(t *testing.T) {

}
