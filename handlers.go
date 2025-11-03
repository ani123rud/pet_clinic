package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"petclinic/data"
	"petclinic/logger"
)

// -------------------- Owners --------------------

func GetOwners(w http.ResponseWriter, r *http.Request) {
	logger.Info("Fetching all owners")
	rows, err := data.ListOwners(DB)
	if err != nil {
		logger.Error("Failed to fetch owners: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	owners := []Owner{}
	for _, ro := range rows {
		owners = append(owners, Owner{ID: ro.ID, Name: ro.Name, Phone: ro.Phone, Address: ro.Address})
	}
	logger.Debug("Retrieved %d owners", len(owners))
	json.NewEncoder(w).Encode(owners)
}

func GetOwnerByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	logger.Debug("Fetching owner with ID: %d", id)
	ro, err := data.GetOwnerByID(DB, id)
	if err != nil {
		logger.Error("Failed to fetch owner with ID %d: %v", id, err)
		http.Error(w, "owner not found", http.StatusNotFound)
		return
	}

	o := Owner{ID: ro.ID, Name: ro.Name, Phone: ro.Phone, Address: ro.Address}
	logger.Debug("Successfully retrieved owner: %+v", o)
	json.NewEncoder(w).Encode(o)
}

func CreateOwner(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating new owner")
	var o Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		logger.Warn("Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	logger.Debug("Processing owner data: %+v", o)
	id, err := data.CreateOwner(DB, data.OwnerInput{
		Name:    o.Name,
		Phone:   o.Phone,
		Address: o.Address,
	})

	if err != nil {
		logger.Error("Failed to create owner: %v", err)
		http.Error(w, "failed to create owner", http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully created owner with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Owner created successfully",
	})
}

// UpdateOwner updates an existing owner
func UpdateOwner(w http.ResponseWriter, r *http.Request) {
	logger.Info("Updating owner")
	
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid owner ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get existing owner
	_, err = data.GetOwnerByID(DB, id)
	if err != nil {
		logger.Error("Owner not found with ID %d: %v", id, err)
		http.Error(w, "owner not found", http.StatusNotFound)
		return
	}

	// Decode request body
	var o Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		logger.Warn("Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Update owner in database
	err = data.UpdateOwner(DB, id, data.OwnerInput{
		Name:    o.Name,
		Phone:   o.Phone,
		Address: o.Address,
	})

	if err != nil {
		logger.Error("Failed to update owner with ID %d: %v", id, err)
		http.Error(w, "failed to update owner", http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully updated owner with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Owner updated successfully",
	})
}

// DeleteOwner deletes an owner by ID
func DeleteOwner(w http.ResponseWriter, r *http.Request) {
	logger.Info("Deleting owner")
	
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid owner ID format for deletion: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Check if owner exists
	_, err = data.GetOwnerByID(DB, id)
	if err != nil {
		logger.Error("Owner not found for deletion with ID %d: %v", id, err)
		http.Error(w, "owner not found", http.StatusNotFound)
		return
	}

	// Delete owner
	err = data.DeleteOwner(DB, id)
	if err != nil {
		logger.Error("Failed to delete owner with ID %d: %v", id, err)
		http.Error(w, "failed to delete owner", http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully deleted owner with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// -------------------- Pets --------------------

// UpdatePet updates an existing pet
func UpdatePet(w http.ResponseWriter, r *http.Request) {
	logger.Info("Updating pet")
	
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid pet ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get existing pet
	_, err = data.GetPetByID(DB, id)
	if err != nil {
		logger.Error("Pet not found with ID %d: %v", id, err)
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	// Decode request body
	var p Pet
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		logger.Warn("Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Update pet in database
	err = data.UpdatePet(DB, id, data.PetInput{
		Name:    p.Name,
		Species: p.Species,
		Breed:   p.Breed,
		Birth:   p.Birth,
		OwnerID: p.OwnerID,
	})

	if err != nil {
		logger.Error("Failed to update pet with ID %d: %v", id, err)
		http.Error(w, "failed to update pet", http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully updated pet with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Pet updated successfully",
	})
}

// DeletePet deletes a pet by ID
func DeletePet(w http.ResponseWriter, r *http.Request) {
	logger.Info("Deleting pet")
	
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid pet ID format for deletion: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Check if pet exists
	_, err = data.GetPetByID(DB, id)
	if err != nil {
		logger.Error("Pet not found for deletion with ID %d: %v", id, err)
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	// Delete pet
	err = data.DeletePet(DB, id)
	if err != nil {
		logger.Error("Failed to delete pet with ID %d: %v", id, err)
		http.Error(w, "failed to delete pet", http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully deleted pet with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// -------------------- Pets --------------------

func GetPets(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Fetching all pets")
	rows, err := data.ListPets(DB)
	if err != nil {
		logger.Error("Failed to fetch pets: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	pets := []Pet{}
	for _, rp := range rows {
		pets = append(pets, Pet{ID: rp.ID, Name: rp.Name, Species: rp.Species, Breed: rp.Breed, Birth: rp.Birth, OwnerID: rp.OwnerID})
	}
	logger.Debug("Retrieved %d pets", len(pets))
	json.NewEncoder(w).Encode(pets)
}

func GetPetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid pet ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	logger.Debug("Fetching pet with ID: %d", id)
	rp, err := data.GetPetByID(DB, id)
	if err != nil {
		logger.Error("Pet not found with ID %d: %v", id, err)
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	p := Pet{ID: rp.ID, Name: rp.Name, Species: rp.Species, Breed: rp.Breed, Birth: rp.Birth, OwnerID: rp.OwnerID}
	logger.Debug("Successfully retrieved pet: %+v", p)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func CreatePet(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating new pet")
	var p Pet
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		logger.Warn("Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	logger.Debug("Processing pet data: %+v", p)
	id, err := data.CreatePet(DB, data.PetInput{
		Name:    p.Name,
		Species: p.Species,
		Breed:   p.Breed,
		Birth:   p.Birth,
		OwnerID: p.OwnerID,
	})

	if err != nil {
		logger.Error("Failed to create pet: %v", err)
		http.Error(w, "failed to create pet", http.StatusInternalServerError)
		return
	}

	p.ID = id
	logger.Info("Successfully created pet with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// -------------------- Visits --------------------

func GetVisits(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Fetching all visits")
	rows, err := data.ListVisits(DB)
	if err != nil {
		logger.Error("Failed to fetch visits: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	visits := []Visit{}
	for _, rv := range rows {
		visits = append(visits, Visit{ID: rv.ID, PetID: rv.PetID, VetID: rv.VetID, Visit: rv.Visit, Desc: rv.Desc})
	}
	logger.Debug("Retrieved %d visits", len(visits))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(visits)
}

func GetVisitByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Warn("Invalid visit ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	logger.Debug("Fetching visit with ID: %d", id)
	rv, err := data.GetVisitByID(DB, id)
	if err != nil {
		logger.Error("Visit not found with ID %d: %v", id, err)
		http.Error(w, "visit not found", http.StatusNotFound)
		return
	}

	v := Visit{ID: rv.ID, PetID: rv.PetID, VetID: rv.VetID, Visit: rv.Visit, Desc: rv.Desc}
	logger.Debug("Successfully retrieved visit: %+v", v)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func CreateVisit(w http.ResponseWriter, r *http.Request) {
	logger.Info("Creating new visit")
	var v Visit
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logger.Warn("Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	logger.Debug("Processing visit data: %+v", v)
	id, err := data.CreateVisit(DB, data.VisitInput{
		PetID: v.PetID,
		VetID: v.VetID,
		Visit: v.Visit,
		Desc:  v.Desc,
	})

	if err != nil {
		logger.Error("Failed to create visit: %v", err)
		http.Error(w, "failed to create visit", http.StatusInternalServerError)
		return
	}

	v.ID = id
	logger.Info("Successfully created visit with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

// UpdateVisit updates an existing visit
func UpdateVisit(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get existing visit
	_, err = data.GetVisitByID(DB, id)
	if err != nil {
		http.Error(w, "visit not found", http.StatusNotFound)
		return
	}

	// Decode request body
	var v Visit
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Update visit in database
	err = data.UpdateVisit(DB, id, data.VisitInput{
		PetID: v.PetID,
		VetID: v.VetID,
		Visit: v.Visit,
		Desc:  v.Desc,
	})

	if err != nil {
		http.Error(w, "failed to update visit", http.StatusInternalServerError)
		return
	}

	// Return updated visit
	updatedVisit, _ := data.GetVisitByID(DB, id)
	v = Visit{
		ID:    updatedVisit.ID,
		PetID: updatedVisit.PetID,
		VetID: updatedVisit.VetID,
		Visit: updatedVisit.Visit,
		Desc:  updatedVisit.Desc,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// DeleteVisit deletes a visit by ID
func DeleteVisit(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Check if visit exists
	_, err = data.GetVisitByID(DB, id)
	if err != nil {
		http.Error(w, "visit not found", http.StatusNotFound)
		return
	}

	// Delete visit
	err = data.DeleteVisit(DB, id)
	if err != nil {
		http.Error(w, "failed to delete visit", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
