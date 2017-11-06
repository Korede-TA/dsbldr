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

func (b *Builder) addFeatureData(featureName string, values []string) error {
	writeStringColumn(b.data, featureName, values)
	return nil
}

func (b *Builder) getFeatureData(featureName string) []string {
	items := make([]string, b.records)
	if _, ok := b.featureMap[featureName]; ok {
		readStringColumn(items, featureName, b.data)
		return items
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
	return feat
}

func (b *Builder) writeRecord(writer csv.Writer, i int) error {
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
	return nil
}

// Save commits the downloaded features to a file
func (b *Builder) Save(writer csv.Writer) error {
	for i := range b.data {
		err := b.writeRecord(writer, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveIf saves records only if they evaluate to true
func (b *Builder) SaveIf(writer csv.Writer, saveCond func(r []string) bool) error {
	for i := range b.data {
		if saveCond(b.data[i]) {
			err := b.writeRecord(writer, i)
			if err != nil {
				return err
			}
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
) (*http.Request, error) {
	headers := b.RequestHeaders

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
			parentValuesMap[j] = b.getFeatureData(j)[i]
		}
		var resolveEndpointError error
		endpoints[i], resolveEndpointError = feature.resolveEndpoint(parentValuesMap)
		if resolveEndpointError != nil {
			return nil, resolveEndpointError
		}
	}

	return endpoints, nil
}

type endpointClient interface {
	Do(req http.Request) (*http.Response, error)
}

// populateFeatureData returns string dumps of responses and an error if any
func (b *Builder) populateFeatureData(feature *Feature, client endpointClient) ([]string, error) {
	responseDumps := make([]string, b.records)
	endpoints, err := b.resolveFeatureEndpoints(feature)
	if err != nil {
		return nil, err
	}

	for i := 0; i < b.records; i++ {
		// endpoints should be the same as the number of
		req, err := b.createRequest(endpoints[i])
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(*req)
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
func (b *Builder) Run(client endpointClient) error {

	for _, feature := range b.featureMap {
		parents, err := feature.getParentNames()
		if err != nil {
			return err
		}

		var populateError error
		go func(feature *Feature) {
			// Block till parent features finish executing
			for _, i := range parents {
				<-b.GetFeature(i).finished
			}

			parsedResponses := make([]string, b.records)
			parsedResponses, populateError = b.populateFeatureData(feature, client)

			output := feature.RunFunc(parsedResponses)
			b.addFeatureData(feature.Name, output)

			feature.finished <- true // Write to feature.finished channel
			close(feature.finished)
		}(feature)

		if populateError != nil {
			return populateError
		}
	}

	// wait for all channels to be written before exiting function
	for _, feature := range b.featureMap {
		<-feature.finished
	}

	return nil
}
