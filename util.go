package dsbldr

import (
	"fmt"
)

// BasicOAuthHeader spits out a basic OAuth Header based on access token
func BasicOAuthHeader(consumerKey, nonce, signature, signatureMethod,
	timestamp, token string) string {
	return fmt.Sprintf(`OAuth oauth_consumer_key="%s",
		oauth_nonce="%s",
		oauth_signature="%s",
		oauth_signature_method="%s",
		oauth_timestamp="%s",
		oauth_token="%s`,
		consumerKey, nonce, signature, signatureMethod, timestamp, token)
}

func writeStringColumn(data *[][]string, columnName string, values []string) {
	var colIndex int
	for i := range (*data)[0] {
		// Find first empty column
		if (*data)[0][i] == "" {
			colIndex = i
			(*data)[0][i] = columnName
			break
		}
	}
	// Add all the values as well (remember that Builder.data is pre-allocated)
	for i := 1; i < len(*data); i++ {
		(*data)[i][colIndex] = values[i-1]
	}
}
