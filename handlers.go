package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"petclinic/data"
	"petclinic/models"
)

// -------------------- Owners --------------------

func GetOwners(w http.ResponseWriter, r *http.Request) {
	owners, err := data.ListOwners()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(owners)
}

func GetOwnerByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	o, err := data.GetOwnerByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(o)
}

func CreateOwner(w http.ResponseWriter, r *http.Request) {
	var o models.Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	created, err := data.CreateOwner(o)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(created)
}

// -------------------- Pets --------------------

func GetPets(w http.ResponseWriter, r *http.Request) {
	pets, err := data.ListPets()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(pets)
}

func GetPetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	p, err := data.GetPetByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func CreatePet(w http.ResponseWriter, r *http.Request) {
	var p models.Pet
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	created, err := data.CreatePet(p)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(created)
}

// -------------------- Visits --------------------

func GetVisits(w http.ResponseWriter, r *http.Request) {
	visits, err := data.ListVisits()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(visits)
}

func GetVisitByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}
	v, err := data.GetVisitByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	json.NewEncoder(w).Encode(v)
}

func CreateVisit(w http.ResponseWriter, r *http.Request) {
	var v models.Visit
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	created, err := data.CreateVisit(v)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(created)
}
