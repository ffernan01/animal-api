package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ffernan01/animal-api/internal/animal"
	"github.com/ffernan01/animal-api/internal/species"
	"github.com/go-chi/chi/v5"
)

func getAnimals(w http.ResponseWriter, r *http.Request) {
	repo := r.Context().Value(animalRepositoryKey).(animal.Repository)
	a, e := repo.GetAll(r.Context())

	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	b, e := json.Marshal(a)
	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	animalID := chi.URLParam(r, "animalID")

	repo := r.Context().Value(animalRepositoryKey).(animal.Repository)
	a, e := repo.Get(r.Context(), animal.ID(animalID))

	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	if a == animal.InvalidOne {
		fmt.Println("animal wasn't found")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	b, e := json.Marshal(a)
	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	animalID := chi.URLParam(r, "animalID")

	repo := r.Context().Value(animalRepositoryKey).(animal.Repository)
	e := repo.Delete(r.Context(), animal.ID(animalID))

	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func postAnimal(w http.ResponseWriter, r *http.Request) {
	var dto struct {
		Name    string `json:"name"`
		Species string `json:"species"`
	}

	creator := r.Context().Value(animalCreatorKey).(animal.Creator)
	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()

	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	e = json.Unmarshal(b, &dto)
	if e != nil {
		fmt.Printf("bad request. %v", e)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	a, e := creator.Create(r.Context(), dto.Name, species.ID(dto.Species))
	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	b, e = json.Marshal(a)
	if e != nil {
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
