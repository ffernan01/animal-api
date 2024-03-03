package animal

import (
	"context"
	"errors"
	"fmt"

	"github.com/ffernan01/animal-api/internal/species"
)

// Creator creates animals, as simple as that
type Creator interface {
	Create(context.Context, string, species.ID) (Animal, error)
}

type defaultCreator struct {
	aniRepo Repository
	speRepo species.Repository
}

// NewDefaultCreator creates and return a Creator that creates animals
func NewDefaultCreator(ar Repository, sr species.Repository) Creator {
	return &defaultCreator{aniRepo: ar, speRepo: sr}
}

// Create creates an animal
func (d *defaultCreator) Create(ctx context.Context, name string, sID species.ID) (Animal, error) {
	a, e := d.aniRepo.FindByName(ctx, name)
	if e != nil {
		return InvalidOne, fmt.Errorf("couldn't check if animal exists. %w", e)
	}

	if a != InvalidOne {
		return InvalidOne, errors.New("there's already an animal with that name")
	}

	s, e := d.speRepo.Get(ctx, sID)
	if e != nil {
		return InvalidOne, fmt.Errorf("couldn't check if species exists. %w", e)
	}

	if s == species.InvalidOne {
		return InvalidOne, errors.New("hey! what kind of species is that?")
	}

	a, e = d.aniRepo.Save(ctx, Animal{Name: name, Species: s})
	if e != nil {
		return InvalidOne, fmt.Errorf("there was an error. %w", e)
	}

	return a, nil
}
