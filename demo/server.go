package main

/*
Demo server that spits out 1000 random records of a certain kind in a response
*/

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

const (
	recordCount    = 1000
	randomIntRange = 10000
)

func randIntStrings(count int) []string {
	output := make([]string, count)
	r := rand.New()
	for i := 0; i < count; i++ {
		output[i] = strconv.Itoa(r.Int31n(randomIntRange))
	}
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	path := url.Parse(r.URL).EscapedPath

	items := make(map[string][]string)
	items["items"] = randIntStrings(recordCount)
	js, err := json.Marshal(items)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func pricesHandler(w http.ResponseWriter, r *http.Request) {
	path := url.Parse(r.URL).EscapedPath

	prices := make(map[string][]string)
	prices["prices"] = randIntStrings(recordCount)
	js, err := json.Marshal(prices)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	categories := []string{"category1", "category2", "category3", "category4"}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/api/items/", itemsHandler)
	http.HandleFunc("/api/items/prices", pricesHandler)
	http.HandleFunc("/api/items/categories", categoryHandler)
	http.ListenAndServe(":8080", nil)
}
