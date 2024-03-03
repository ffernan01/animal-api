package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ffernan01/animal-api/internal/species"
)

func getSpecies(w http.ResponseWriter, r *http.Request) {
	repo := r.Context().Value(speciesRepositoryKey).(species.Repository)
	s, e := repo.GetAll(r.Context())

	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	b, e := json.Marshal(s)
	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
