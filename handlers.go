package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"petclinic/data"
)

// -------------------- Owners --------------------

func GetOwners(w http.ResponseWriter, r *http.Request) {
	rows, err := data.ListOwners(DB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	owners := []Owner{}
	for _, ro := range rows {
		owners = append(owners, Owner{ID: ro.ID, Name: ro.Name, Phone: ro.Phone, Address: ro.Address})
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
	ro, err := data.GetOwnerByID(DB, id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	o := Owner{ID: ro.ID, Name: ro.Name, Phone: ro.Phone, Address: ro.Address}
	json.NewEncoder(w).Encode(o)
}

func CreateOwner(w http.ResponseWriter, r *http.Request) {
	var o Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	id, err := data.CreateOwner(DB, data.OwnerInput{Name: o.Name, Phone: o.Phone, Address: o.Address})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	o.ID = id
	json.NewEncoder(w).Encode(o)
}

// -------------------- Pets --------------------

func GetPets(w http.ResponseWriter, r *http.Request) {
	rows, err := data.ListPets(DB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	pets := []Pet{}
	for _, rp := range rows {
		pets = append(pets, Pet{ID: rp.ID, Name: rp.Name, Species: rp.Species, Breed: rp.Breed, Birth: rp.Birth, OwnerID: rp.OwnerID})
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
	rp, err := data.GetPetByID(DB, id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	p := Pet{ID: rp.ID, Name: rp.Name, Species: rp.Species, Breed: rp.Breed, Birth: rp.Birth, OwnerID: rp.OwnerID}
	json.NewEncoder(w).Encode(p)
}

func CreatePet(w http.ResponseWriter, r *http.Request) {
	var p Pet
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	id, err := data.CreatePet(DB, data.PetInput{Name: p.Name, Species: p.Species, Breed: p.Breed, Birth: p.Birth, OwnerID: p.OwnerID})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	p.ID = id
	json.NewEncoder(w).Encode(p)
}

// -------------------- Visits --------------------

func GetVisits(w http.ResponseWriter, r *http.Request) {
	rows, err := data.ListVisits(DB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	visits := []Visit{}
	for _, rv := range rows {
		visits = append(visits, Visit{ID: rv.ID, PetID: rv.PetID, VetID: rv.VetID, Visit: rv.Visit, Desc: rv.Desc})
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
	rv, err := data.GetVisitByID(DB, id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	v := Visit{ID: rv.ID, PetID: rv.PetID, VetID: rv.VetID, Visit: rv.Visit, Desc: rv.Desc}
	json.NewEncoder(w).Encode(v)
}

func CreateVisit(w http.ResponseWriter, r *http.Request) {
	var v Visit
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	id, err := data.CreateVisit(DB, data.VisitInput{PetID: v.PetID, VetID: v.VetID, Visit: v.Visit, Desc: v.Desc})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	v.ID = id
	json.NewEncoder(w).Encode(v)
}
