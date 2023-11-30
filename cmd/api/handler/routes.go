package handler

import (
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

func Configure(r chi.Router) {
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
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Heartbeat("/ping"))

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
}
