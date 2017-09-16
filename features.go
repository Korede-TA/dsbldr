package dsbldr

// BaseFeature Interface
type BaseFeature interface {
	Run()
	Deserialize()
}

// PrimaryFeature in the dataset, on which all other features are based on
// Implements BaseFeature
type PrimaryFeature struct {
	Name        string
	URLEndpoint string // Enpoint on social API that the feature aggregates data from
	data        string // JSON string containing data
}

// Run retrieves the feature via JSON GET requests to the API
func (pf PrimaryFeature) Run() {

}

// Deserialize returns a map that represents the JSON deserialized data field
func (pf PrimaryFeature) Deserialize() map[string]interface{} {
	return nil
}

// Feature is a Regular dataset feature that is based off the primary field
// Implements BaseFeature
type Feature struct {
	Name        string
	URLEndpoint string // Enpoint on social API that the feature aggregates data from
	data        string // JSON string containing data
}

// Run retrieves the feature via JSON GET requests to the API
func (f Feature) Run() {

}

// Deserialize returns a map that represents the JSON deserialized data field
func (f Feature) Deserialize() map[string]interface{} {
	return nil
}
