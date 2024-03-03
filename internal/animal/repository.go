package animal

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Repository is an animal repository abstraction
type Repository interface {
	Get(context.Context, ID) (Animal, error)
	GetAll(context.Context) ([]Animal, error)
	FindByName(context.Context, string) (Animal, error)
	Save(context.Context, Animal) (Animal, error)
	Delete(context.Context, ID) error
}

type inMemRepository struct {
	storage sync.Map
}

// NewInMemRepository creates a Repository that stores animals in memory
func NewInMemRepository() Repository {
	return &inMemRepository{}
}

// GetAll gets all stored animals
func (i *inMemRepository) GetAll(context.Context) ([]Animal, error) {
	animals := make([]Animal, 0, 10)
	var aniVal Animal

	i.storage.Range(func(key, value any) bool {
		aniVal = value.(Animal)
		animals = append(animals, aniVal)
		return true
	})

	return animals, nil
}

// Get returns an animal if it's already stored
func (i *inMemRepository) Get(ctx context.Context, id ID) (Animal, error) {
	v, ok := i.storage.Load(id)
	if !ok {
		return InvalidOne, nil
	}

	return v.(Animal), nil
}

// FindByName finds an animal by its name
func (i *inMemRepository) FindByName(ctx context.Context, name string) (Animal, error) {
	v := InvalidOne
	var aniVal Animal

	i.storage.Range(func(key, value any) bool {
		aniVal = value.(Animal)
		if strings.ToLower(aniVal.Name) == strings.ToLower(name) {
			v = aniVal
			return false
		} else {
			return true
		}
	})
	return v, nil
}

// Save saves an animal in a inmem storage
func (i *inMemRepository) Save(ctx context.Context, animal Animal) (Animal, error) {
	animal.ID = ID(uuid.New().String())
	i.storage.Store(animal.ID, animal)
	return animal, nil
}

// Delete deletes an animal
func (i *inMemRepository) Delete(ctx context.Context, id ID) error {
	i.storage.Delete(id)
	return nil
}
