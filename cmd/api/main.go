package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ffernan01/animal-api/internal/animal"
	"github.com/ffernan01/animal-api/internal/species"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	animalRepositoryKey  = "animal-repository"
	animalCreatorKey     = "animal-creator"
	speciesRepositoryKey = "species-repository"
)

func main() {
	r := chi.NewRouter()
	// initialize deps
	ar := animal.NewInMemRepository()
	sr := species.NewInMemRepository()
	ac := animal.NewDefaultCreator(ar, sr)

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// ioc
	r.Use(middleware.WithValue(animalRepositoryKey, ar))
	r.Use(middleware.WithValue(speciesRepositoryKey, sr))
	r.Use(middleware.WithValue(animalCreatorKey, ac))

	// routes
	r.Route("/animals", func(r chi.Router) {
		r.Get("/", getAnimals)
		r.Post("/", postAnimal)

		r.Route("/{animalID}", func(r chi.Router) {
			r.Get("/", getAnimal)
			r.Delete("/", deleteAnimal)
		})
	})

	r.Route("/species", func(r chi.Router) {
		r.Get("/", getSpecies)
	})

	http.ListenAndServe(":8080", r)
}

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
		fmt.Printf("an error occurred. %v", e)
		http.Error(w, "an error occurred", http.StatusBadRequest)
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
