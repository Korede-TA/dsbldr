package dsbldr

// RunFunc holds the computation that processes the API responses to features
// is sent a JSON string and returns
type RunFunc func(data string) []interface{}

// Feature in the dataset, on which all other features are based on
// Implements BaseFeature
type Feature struct {
	Name     string
	Endpoint string // API Endpoint
	RunFunc  RunFunc
}

// NewFeature creates new Feature with defaults
func NewFeature() *Feature {
	return &Feature{}
}
