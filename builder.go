package dsbldr

import (
	"encoding/csv"
	"io/ioutil"
	"net/http"
)

// Builder is main type for this tool.
type Builder struct {
	BaseURL        string
	RequestHeaders map[string]string // Custom Request Headers including auth
	featureMap     map[string]*Feature
	data           [][]string // Strings of Data to be read in to CSV
	records        int        // Number of records to be retrieved for dataset
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
		records:        recordCount,
	}
}

func (b *Builder) addDataFeature(featureName string, values []string) error {
	writeStringColumn(&b.data, featureName, values)
	return nil
}

func (b *Builder) getDataFeature(featureName string) []string {
	items := make([]string, b.records-1) // all items excluding header row
	var colIndex int

	for i := range b.data[0] {
		// Find first column with same header
		if b.data[0][i] == featureName {
			colIndex = i
			break
		}
	}

	// Add all the values as well (remember that Builder.data is pre-allocated)
	for i := 1; i < len(b.data); i++ {
		items[i-1] = b.data[i][colIndex]
	}

	return items
}

// GetFeature returns a feature in the detaset based on it's name
func (b *Builder) GetFeature(name string) *Feature {
	var feat *Feature
	if val, ok := b.featureMap[name]; ok {
		feat = val
		return feat
	}
	return nil
}

// Save commits the downloaded features to a file
func (b *Builder) Save(writer csv.Writer) error {
	for i := range b.data {
		var record []string
		for index, j := range b.data[i] {
			// if data header (feature name) has noSave == false, else don't write
			if !b.featureMap[b.data[0][index]].noSave {
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
		// TODO: Return error if feature with same name has been added
	}
	// Increase size of data matrix if feature map is larger than initially allocated
	if len(b.featureMap) > len(b.data[0]) {
		for i := range b.data {
			b.data[i] = append(b.data[i], "")
		}
	}
}

func (b *Builder) createRequest(
	endpoint string,
	headers map[string]string,
) (*http.Request, error) {

	req, err := http.NewRequest("GET", b.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req, nil
}

// resolve endpoint templates to actual values for endpoints based on parent features
func (b *Builder) resolveFeatureEndpoints(feature *Feature) ([]string, error) {
	endpoints := make([]string, b.records)
	parents, err := feature.getParentNames()
	if err != nil {
		return nil, err
	}

	for i := 0; i < b.records; i++ {
		parentValuesMap := make(map[string]string)
		for _, j := range parents {
			parentValuesMap[j] = b.getDataFeature(j)[i]
		}
		var resolveEndpointError error
		endpoints[i], resolveEndpointError = feature.resolveEndpoint(parentValuesMap)
		if resolveEndpointError != nil {
			return nil, resolveEndpointError
		}
	}

	return endpoints, nil
}

// retrieveFeature returns string dumps of responses and an error if any
func (b *Builder) retrieveFeature(feature *Feature) ([]string, error) {
	client := &http.Client{}
	responses := make([]*http.Response, b.records)
	responseDumps := make([]string, b.records)
	endpoints, err := b.resolveFeatureEndpoints(feature)
	if err != nil {
		return nil, err
	}

	for i := 0; i < b.records; i++ {
		// endpoints should be the same as the number of
		req, err := b.createRequest(endpoints[i], b.RequestHeaders)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode == 200 { // OK
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			responseDumps[i] = string(bodyBytes)
		}
	}

	return responseDumps, nil
}

// Run Builder to aggregate all features and manage concurrent operations
func (b *Builder) Run() error {
	// flag.Parse()

	// Set noSave to true for any features listed in the call
	// for _, i := range ignoredFeatures {
	// 	b.GetFeature(i).noSave = true
	// }

	for _, feature := range b.featureMap {
		parents, err := feature.getParentNames()
		if err != nil {
			return err
		}

		var retrieveError error
		go func(feature *Feature) {
			for _, i := range parents {
				// Loop through parent features and wait till
				// they finish excecuting
				<-b.GetFeature(i).finished
			}

			parsedResponses := make([]string, b.records)
			parsedResponses, retrieveError = b.retrieveFeature(feature)

			output := feature.RunFunc(parsedResponses)
			b.addDataFeature(feature.Name, output)

			feature.finished <- true // Mark feature as true
		}(feature)

		if retrieveError != nil {
			return retrieveError
		}

	}
	return nil
}
