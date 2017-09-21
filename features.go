package dsbldr

// RunFunc holds the computation that processes the API responses to features
// is sent a JSON string of the response ??as well as a map of data from the features parent features??
// Basically what you do with the run function is take in a string of
// serialized API data (could be in JSON or XML), do parsing on your
// own or using utility functions. You do whatever computations you want and
// then spit it back as an array of strings to read to CSV or JSON
type RunFunc func(response string) []string // parents map[string]string

// const (
// 	EmptyFeature = iota
// 	FullFeature
// )

// Feature in the dataset, on which all other features are based on
type Feature struct {
	Name     string
	Endpoint string // API Endpoint
	RunFunc  RunFunc
	// status   int
}

// NewFeature creates new Feature with defaults
func NewFeature() *Feature {
	return &Feature{}
}
