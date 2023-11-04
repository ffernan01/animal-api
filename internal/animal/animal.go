package animal

import "github.com/ffernan01/animal-api/internal/species"

var InvalidOne Animal = Animal{ID: "INVALID", Name: "INVALID", Species: species.InvalidOne}

// ID animal identifier
type ID string

// Animal is a reprentative entity of the animal kingdom
type Animal struct {
	ID      ID
	Name    string
	Species species.Species
}
