package species

import (
	"context"
	"sync"
)

var (
	InvalidOne Species = Species{
		ID:   "INVALID",
		Name: "INVALID",
	}
)

// ID species identifier
type ID string

// Species is an animal species, like a family
type Species struct {
	ID   ID
	Name string
}

// Repository is a species repository abstraction
type Repository interface {
	Get(context.Context, ID) (Species, error)
	GetAll(context.Context) ([]Species, error)
}

type inMemRepository struct {
	storage sync.Map
}

// NewInMemRepository creates a Repository that stores species in memory
func NewInMemRepository() Repository {
	i := &inMemRepository{}

	s := Species{ID: "ave", Name: "Aves"}
	i.storage.Store(s.ID, s)
	s = Species{ID: "paq", Name: "Paquidermos"}
	i.storage.Store(s.ID, s)
	s = Species{ID: "can", Name: "Canidos"}
	i.storage.Store(s.ID, s)
	s = Species{ID: "fel", Name: "Felinos"}
	i.storage.Store(s.ID, s)

	return i
}

// Get gets an species by its ID
func (i *inMemRepository) Get(ctx context.Context, id ID) (Species, error) {
	v, ok := i.storage.Load(id)
	if !ok {
		return InvalidOne, nil
	}

	return v.(Species), nil
}

// GetAll gets all stored species
func (i *inMemRepository) GetAll(context.Context) ([]Species, error) {
	species := make([]Species, 0, 10)
	var speVal Species

	i.storage.Range(func(key, value any) bool {
		speVal = value.(Species)
		species = append(species, speVal)
		return true
	})

	return species, nil
}
