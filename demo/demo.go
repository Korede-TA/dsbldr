package main

import (
	d "dsbldr"
	"encoding/json"
	"fmt"
)

builder d.Builder := d.Builder{
	BaseURL: "localhost:8080",
	RequestHeaders: map[string]string{
		"Authorization": BasicOAuthHeader(
			"OAUTH_CONSUMER_KEY",
			"OAUTH_NONCE",
			"OAUTH_SIGNATURE",
			"OAUTH_SIGNATURE_METHOD", "OAUTH_TIMESTAMP",
			"OAUTH_TOKEN",
		),
	},
}

builder.AddFeatures(
	&d.Feature{
		Name: "item_ids",
		Endpoint: "/items/",
		RunFunc: func(response []string) []string {
			responseMap = (make[map]string)
			json.Unmarshal(response, &responseMap)
		},
	},
	&d.Feature{
		Name: "item_prices",
		Endpoint: "/items/prices/{{item_ids}}/",
		RunFunc: func(response []string) []string {
			// blah blah
		},
	},
	&d.Feature{
		Name: "item_category",
		Endpoint: "/items/category/{{item_ids}}/",
		RunFunc: func(response []string) []string {
			// blah blah
		},
	},
)

//func main() {
//}
