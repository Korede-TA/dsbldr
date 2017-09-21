package dsbldr

import (
	"io"
)

// Builder is main type for this tool.
type Builder struct {
	BaseURL  string // Base API url to use
	Primary  Feature
	Features []Feature // Normal features based on primary
}

// func (b *Builder) buildFeatureTree() {
// 	// Since datasets might consist of a number of interdependent features
// 	// The tool will be tasked with building a tree of the features

// }

// Init initializes the Builder
func (b *Builder) Init() {

}

// Run Builder to aggregate all features and manage concurrent operations
func (b *Builder) Run() error {
	return nil
}

// Save commits the
func (b *Builder) Save(io.Writer) error {
	return nil
}

// AddFeature adds a Feature struct to the "Features" Field on Builder
func (b *Builder) AddFeature(feature Feature) {
	b.Features = append(b.Features, feature)
}

// GetFeature returns a feature in the detaset based on it's name
func (b *Builder) GetFeature(name string) *Feature {
	for _, i := range b.Features {
		if i.Name == name {
			return &i
		}
	}
	return nil
}
