package print

import "github.com/goliatone/go-masker"

var PrintMasker Masker = defaultMasker{masker: *masker.Default}

// Masker defines the interface for masking data
type Masker interface {
	// Mask takes any type T and returns a masked version of the same type
	// along with any error that occurred during masking
	Mask(target any) (ret any, err error)
}

// NoOpMasker implements Masker but doesn't modify the data
type defaultMasker struct {
	masker masker.Masker
}

func (d defaultMasker) Mask(target any) (ret any, err error) {
	return d.masker.Mask(target)
}
