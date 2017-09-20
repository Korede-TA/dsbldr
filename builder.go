package dsbldr

import (
	"io"
	"net/http"
)

// Builder is main type for this tool.
type Builder struct {
	BaseURL        string
	RequestHeaders map[string]string // Manually written Request Header (including auth)
	FeatureMap     map[string]*Feature
	data           [][]string // Strings of Data to be read in to CSV
	request        http.Request
}

// NewBuilder creates new Builder struct
func NewBuilder(features, recordCount int) *Builder {
	preallocatedData := make([][]string, features)
	for i := range preallocatedData {
		preallocatedData[i] = make([]string, recordCount)
	}
	return &Builder{
		data: preallocatedData,
	}
}

func (b *Builder) createClientAndRequest(endpoint string, headers map[string]string) (*http.Client, *http.Request, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.BaseURL+endpoint, nil)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return client, req, err
}

// Init initializes the Builder
func (b *Builder) Init() {
}

// Run Builder to aggregate all features and manage concurrent operations
func (b *Builder) Run() error {
	for _, feature := range b.FeatureMap {
		client, req, err := b.createClientAndRequest(feature.Endpoint, b.RequestHeaders)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
	}
	return nil
}

// Save commits the downloaded features to a file
// as specified by the Writer interface (has to implement WriteAll)
func (b *Builder) Save(writer io.Writer) error {
	err := writer.WriteAll(b.data)
	return err
}

// AddFeatures adds a Feature struct to the "Features" Field on Builder
func (b *Builder) AddFeatures(features ...*Feature) {
	for _, feature := range features {
		b.FeatureMap[feature.Name] = feature
	}
}

// GetFeature returns a feature in the detaset based on it's name
func (b *Builder) GetFeature(name string) *Feature {
	if val, ok := b.FeatureMap[name]; ok {
		return val
	}
	return nil
}
