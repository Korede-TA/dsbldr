package dsbldr

import (
	"fmt"
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
