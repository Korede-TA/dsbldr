package dsbldr

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBasicOAuthHeader(t *testing.T) {
	consumerKey := "consumerKey"
	nonce := "nonce"
	signature := "signature"
	signatureMethod := "signatureMethod"
	timestamp := "timestamp"
	token := "token"

	want := fmt.Sprintf(`OAuth oauth_consumer_key="%s",
		oauth_nonce="%s",
		oauth_signature="%s",
		oauth_signature_method="%s",
		oauth_timestamp="%s",
		oauth_token="%s`,
		consumerKey, nonce, signature, signatureMethod, timestamp, token)

	got := BasicOAuthHeader(consumerKey, nonce, signature, signatureMethod,
		timestamp, token)

	if got != want {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestWriteStringColumn(t *testing.T) {
	data := [][]string{
		[]string{"a", "b", ""},
		[]string{"a", "b", ""},
		[]string{"a", "b", ""},
		[]string{"a", "b", ""},
	}
	colName := "c"
	values := []string{"c", "c", "c"}

	writeStringColumn(data, colName, values)
	want := [][]string{
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
	}

	if got := data; !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestWriteStringColumnOverwrite(t *testing.T) {
	data := [][]string{
		[]string{"a", "b", "c"},
		[]string{"a", "b", "0"},
		[]string{"a", "b", "0"},
		[]string{"a", "b", "0"},
	}
	colName := "c"
	values := []string{"c", "c", "c"}

	writeStringColumn(data, colName, values)
	want := [][]string{
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
	}

	if got := data; !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestReadStringColumn(t *testing.T) {
	data := [][]string{
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
	}
	colName := "c"
	values := make([]string, 3)

	want := []string{"c", "c", "c"}
	readStringColumn(values, colName, data)

	if got := values; !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}

func TestReadStringColumnWithNoHeader(t *testing.T) {
	data := [][]string{
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
		[]string{"a", "b", "c"},
	}
	colName := "d"
	values := make([]string, 3)

	want := []string{"", "", ""}
	readStringColumn(values, colName, data)

	if got := values; !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %v\n want: %v\n ", got, want)
	}
}
