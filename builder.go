package dsbldr

// Builder is main type for this tool.
type Builder struct {
	BaseURL  string // Base social API url to use
	Primary  PrimaryFeature
	Features []BaseFeature // Normal features whose retrieval is based on the primnary feature
}

// Run Builder to aggregate all features and manage concurrency
func (b *Builder) Run() {

}

// GetFeature returns a specific feature in the detaset based on it's name
func (b *Builder) GetFeature(name string) {

}
