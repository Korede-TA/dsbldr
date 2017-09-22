package dsbldr

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Builder is main type for this tool.
type Builder struct {
	BaseURL        string
	RequestHeaders map[string]string // Manually written Request Header (including auth)
	featureMap     map[string]*Feature
	data           [][]string // Strings of Data to be read in to CSV
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
	// First row is table headers
	var colIndex int
	for i := range b.data[0] {
		// Find first empty column
		if b.data[0][i] == "" {
			colIndex = i
			b.data[0][i] = featureName
			break
		}
	}
	// Add all the values as well (remember that Builder.data is pre-allocated)
	for i := 1; i < len(b.data); i++ {
		b.data[i][colIndex] = values[i-1]
	}
	return nil
}

func (b *Builder) createClientAndRequest(endpoint string, headers map[string]string) (*http.Client, *http.Request, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.BaseURL+endpoint, nil)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return client, req, err
}

// Run Builder to aggregate all features and manage concurrent operations
func (b *Builder) Run() error {
	for _, feature := range b.featureMap {
		client, req, err := b.createClientAndRequest(feature.Endpoint, b.RequestHeaders)
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
	// err := writer.WriteAll(b.data)
	return nil
}

// AddFeatures adds a Feature struct to the "Features" Field on Builder
func (b *Builder) AddFeatures(features ...*Feature) {
	for _, feature := range features {
		b.featureMap[feature.Name] = feature
	}
}

// GetFeature returns a feature in the detaset based on it's name
func (b *Builder) GetFeature(name string) *Feature {
	if val, ok := b.featureMap[name]; ok {
		return val
	}
	return nil
}
