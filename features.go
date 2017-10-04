package dsbldr

import (
	"fmt"
	"regexp"
)

// RunFunc holds the computation that processes the API responses to features
// is sent an array of JSON strings as the responses ??as well as a map of data from the features parent features??
// Basically what you do with the run function is take in a string of
// serialized API data (could be in JSON or XML), do parsing on your
// own or using utility functions. You do whatever computations you want and
// then spit it back as an array of strings to read to CSV or JSON
type RunFunc func(responses []string) []string // parents map[string]string

// Structs representing RetreiveType
// SingleRetrieve Features only require one request to create the JSON Dump
// that's passed to the RunFunc
// Repeated Retrieve Features require one request per value-set of
// of parent features that are concatenated into a JSON array and then passed
// to the Features RunFunc
// Almost as a given, all dependent features will be of RepeatedRetrieve per
// value sets of their parent features
const (
	SingleRetrieve = iota
	RepeatedRetrieve
)

// Feature in the dataset, on which all other features are based on
type Feature struct {
	Name         string
	Endpoint     string  // API Endpoint
	RunFunc      RunFunc // function that performs ad-hoc computation
	RetrieveType int     // Determines if multiple or single requests are made to the api
	noSave       bool
	finished     chan bool // if feature is finished
}

// NewFeature creates new Feature with defaults
func NewFeature() *Feature {
	return &Feature{
		noSave:   false,
		finished: make(chan bool, 1),
	}
}

// getParentNames returns strings of
func (f *Feature) getParentNames() ([]string, error) {
	// parse through using regexp
	re, err := regexp.Compile(`{{\w+}}`)
	if err != nil {
		return nil, err
	}
	val := re.FindAllString(f.Endpoint, -1)
	for i := range val {
		val[i] = val[i][2 : len(val[i])-2] // trim off parentheses
	}
	return val, nil
}

// resolveEndpoint returns the Endpoint but with the actual parent values for that record inserted
func (f *Feature) resolveEndpoint(parentsValues map[string]string) (string, error) {
	toReplace := f.Endpoint

	for k, v := range parentsValues {
		re, err := regexp.Compile(fmt.Sprintf("\\{\\{%s\\}\\}", k))
		if err != nil {
			return "", err
		}
		// overwrite toReplace to (kind of) recursively do any replaces
		toReplace = re.ReplaceAllLiteralString(toReplace, v)
	}

	return toReplace, nil
}
