package main

import (
	"log"
	"net/http"
	"petclinic/data"
)

func main() {
	if err := data.Init(); err != nil {
		log.Fatal(err)
	}

	// Owners
	http.HandleFunc("/owners", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetOwners(w, r)
		} else if r.Method == http.MethodPost {
			CreateOwner(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	http.HandleFunc("/owners/id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetOwnerByID(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	// Pets
	http.HandleFunc("/pets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetPets(w, r)
		} else if r.Method == http.MethodPost {
			CreatePet(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	http.HandleFunc("/pets/id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetPetByID(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	// Visits
	http.HandleFunc("/visits", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetVisits(w, r)
		} else if r.Method == http.MethodPost {
			CreateVisit(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	http.HandleFunc("/visits/id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetVisitByID(w, r)
		} else {
			http.Error(w, "Method not allowed", 405)
		}
	})

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
