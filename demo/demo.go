package main

import (
	d "dsbldr"
	"encoding/json"
	"fmt"
)

builder := d.Builder{
	BaseURL: "localhost:8080",
	RequestHeaders: map[string]string{
		"Authorization": BasicOAuthHeader(
			"oauth_consumer_key",
			"oauth_nonce",
			"oauth_signature",
			"oauth_signature_method", "oauth_timestamp",
			"oauth_token",
		)
	}
}

builder.AddFeatures(
	&d.Feature{
		Name: "item_ids",
		Endpoint: "/items/",
		RunFunc: func(response string) []string {
			responseMap = (make[map]string)
			json.Unmarshal(response, &responseMap)
		}
	},
	&d.Feature{
		Name: "item_prices",
		Endpoint: "/items/prices/{{item_ids}}/",
		RunFunc: func(response string) []string {
			// blah blah
		}
	},
	&d.Feature{
		Name: "item_category",
		Endpoint: "/items/category/{{item_ids}}/",
		RunFunc: func(response string) []string {
			// blah blah
		}
	},
)

func main() {
	fmt.Print(err)
}