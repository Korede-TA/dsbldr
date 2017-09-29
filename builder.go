package dsbldr

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// FeatureList is a custom flag type
type FeatureList []string

// String return value of custom flag type
func (fl *FeatureList) String() string {
	return fmt.Sprint(*fl)
}

// Set value of custom flag type
func (fl *FeatureList) Set(value string) error {
	features := strings.Split(value, " ")
	for _, i := range features {
		*fl = append(*fl, i)
	}
	return nil
}

var featureList FeatureList

func init() {
	flag.Var(
		&featureList,
		"features",
		"comma-separated list of features. Defaults to all features",
	)
}

// Builder is main type for this tool.
type Builder struct {
	BaseURL        string
	RequestHeaders map[string]string // Custom Request Headers including auth
	featureMap     map[string]*Feature
	data           [][]string // Strings of Data to be read in to CSV
	saveFeatures   map[string]bool
}

// NewBuilder creates new Builder struct
func NewBuilder(featureCount, recordCount int) *Builder {
	// Add extra row for header
	preallocatedData := make([][]string, recordCount+1)
	for i := range preallocatedData {
		preallocatedData[i] = make([]string, featureCount)
	}
	return &Builder{
		RequestHeaders: make(map[string]string),
		featureMap:     make(map[string]*Feature),
		data:           preallocatedData,
	}
}

func (b *Builder) addDataFeature(featureName string, values []string) error {
	writeStringColumn(&b.data, featureName, values)
	return nil
}

func (b *Builder) createClientAndRequest(
	endpoint string,
	headers map[string]string,
) (*http.Client, *http.Request, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.BaseURL+endpoint, nil)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return client, req, err
}

// Run Builder to aggregate all features and manage concurrent operations
func (b *Builder) Run() error {
	flag.Parse()

	// Set noSave to true for any features listed in the call
	for _, i := range ignoredFeatures {
		b.GetFeature(i).noSave = true
	}

	for _, feature := range b.featureMap {
		go func(){
			// Run operations concurrently here 
			// (use channels representing download the status 
			// of each record to get values)
			// Requests should be cached since several requests 
			// will have many duplicates (as some will share the same response)

			// TODO: 
				// - Figure out caching stuff within createClientAndRequest
				// -- (Maybe just do the caching manually)
				// - Use channels to make certain download operations 
				//   wait till the dependent values are loaded
		}()
		client, req, err := b.createClientAndRequest(
			feature.Endpoint, b.RequestHeaders)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		parsedResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		output := feature.RunFunc(string(parsedResponse))
		b.addDataFeature(feature.Name, output)
	}
	return nil
}

// Save commits the downloaded features to a file
// as specified by the Writer interface (has to implement WriteAll)
func (b *Builder) Save(writer io.Writer) error {
	for _, i := range b.data {
		var record []string
		for index, j := range b.data[i] {
			// if data header (feature name) has noSave == false, else don't write
			if !b.featureMap[b.data[0][index]].noSave {
				// find out if Golang append acts similarly to lists in pytohn
				// if not, implement special version in util.go
				record = append(record, j)
			}
		}
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddFeatures adds a Feature struct to the "Features" Field on Builder
func (b *Builder) AddFeatures(features ...*Feature) {
	for _, feature := range features {
		feature.noSave = false
		b.featureMap[feature.Name] = feature
		b.saveFeatures[feature.Name] = true
	}
	// Increase size of data if feature map is larger than initially allocated
	if len(b.featureMap) > len(b.data[0]) {
		for i := range b.data {
			b.data[i] = append(b.data[i], "")
		}
	}
}

// GetFeature returns a feature in the detaset based on it's name
func (b *Builder) GetFeature(name string) *Feature {
	var feat *Feature
	if val, ok := b.featureMap[name]; ok {
		feat = val
		return feat
	}
	return feat
}
