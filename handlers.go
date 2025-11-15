package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"petclinic/data"
	"petclinic/logger"
)

// -------------------- Owners --------------------

func GetOwners(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Fetching all owners")
	rows, err := data.ListOwners(DB)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to fetch owners: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	owners := []Owner{}
	for _, ro := range rows {
		owners = append(owners, Owner{ID: ro.ID, Name: ro.Name, Phone: ro.Phone, Address: ro.Address})
	}
	logger.DebugCtx(r.Context(), "Retrieved %d owners", len(owners))
	json.NewEncoder(w).Encode(owners)
}

func GetOwnerByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Fetching owner with ID: %d", id)
	ro, err := data.GetOwnerByID(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to fetch owner with ID %d: %v", id, err)
		http.Error(w, "owner not found", http.StatusNotFound)
		return
	}

	o := Owner{ID: ro.ID, Name: ro.Name, Phone: ro.Phone, Address: ro.Address}
	logger.DebugCtx(r.Context(), "Successfully retrieved owner: %+v", o)
	json.NewEncoder(w).Encode(o)
}

func CreateOwner(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Creating new owner")
	var o Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		logger.WarnCtx(r.Context(), "Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Processing owner data: %+v", o)
	id, err := data.CreateOwner(DB, data.OwnerInput{
		Name:    o.Name,
		Phone:   o.Phone,
		Address: o.Address,
	})

	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to create owner: %v", err)
		http.Error(w, "failed to create owner", http.StatusInternalServerError)
		return
	}

	logger.InfoCtx(r.Context(), "Successfully created owner with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Owner created successfully",
	})
}

// UpdateOwner updates an existing owner
func UpdateOwner(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Updating owner")
	
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid owner ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get existing owner
	_, err = data.GetOwnerByID(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Owner not found with ID %d: %v", id, err)
		http.Error(w, "owner not found", http.StatusNotFound)
		return
	}

	// Decode request body
	var o Owner
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		logger.WarnCtx(r.Context(), "Invalid JSON in request: %v", err)
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
		logger.ErrorCtx(r.Context(), "Failed to update owner with ID %d: %v", id, err)
		http.Error(w, "failed to update owner", http.StatusInternalServerError)
		return
	}

	logger.InfoCtx(r.Context(), "Successfully updated owner with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Owner updated successfully",
	})
}

// DeleteOwner deletes an owner by ID
func DeleteOwner(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Deleting owner")
	
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid owner ID format for deletion: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Check if owner exists
	_, err = data.GetOwnerByID(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Owner not found for deletion with ID %d: %v", id, err)
		http.Error(w, "owner not found", http.StatusNotFound)
		return
	}

	// Delete owner
	err = data.DeleteOwner(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to delete owner with ID %d: %v", id, err)
		http.Error(w, "failed to delete owner", http.StatusInternalServerError)
		return
	}

	logger.InfoCtx(r.Context(), "Successfully deleted owner with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
    logger.InfoCtx(r.Context(), "Uploading file")
    if err := r.ParseMultipartForm(32 << 20); err != nil {
        logger.WarnCtx(r.Context(), "Failed to parse multipart form: %v", err)
        http.Error(w, "invalid multipart form", http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        logger.WarnCtx(r.Context(), "Missing file in form: %v", err)
        http.Error(w, "missing file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    filename := filepath.Base(header.Filename)
    if filename == "" {
        http.Error(w, "invalid filename", http.StatusBadRequest)
        return
    }

    if err := os.MkdirAll("uploads", 0755); err != nil {
        logger.ErrorCtx(r.Context(), "Failed to ensure uploads directory: %v", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    dstPath := filepath.Join("uploads", filename)
    dst, err := os.Create(dstPath)
    if err != nil {
        logger.ErrorCtx(r.Context(), "Failed to create destination file: %v", err)
        http.Error(w, "failed to save file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    n, err := io.Copy(dst, file)
    if err != nil {
        logger.ErrorCtx(r.Context(), "Failed to write file: %v", err)
        http.Error(w, "failed to save file", http.StatusInternalServerError)
        return
    }

    logger.InfoCtx(r.Context(), "Uploaded file: %s (%d bytes)", filename, n)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "name":   filename,
        "size":   n,
        "message": "File uploaded successfully",
    })
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        http.Error(w, "missing name parameter", http.StatusBadRequest)
        return
    }

    filename := filepath.Base(name)
    path := filepath.Join("uploads", filename)

    if _, err := os.Stat(path); err != nil {
        if os.IsNotExist(err) {
            logger.WarnCtx(r.Context(), "File not found: %s", filename)
            http.Error(w, "file not found", http.StatusNotFound)
            return
        }
        logger.ErrorCtx(r.Context(), "Failed to access file: %v", err)
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    logger.InfoCtx(r.Context(), "Downloading file: %s", filename)
    http.ServeFile(w, r, path)
}

// -------------------- Pets --------------------

// UpdatePet updates an existing pet
func UpdatePet(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Updating pet")
	
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid pet ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get existing pet
	_, err = data.GetPetByID(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Pet not found with ID %d: %v", id, err)
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	// Decode request body
	var p Pet
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		logger.WarnCtx(r.Context(), "Invalid JSON in request: %v", err)
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
		logger.ErrorCtx(r.Context(), "Failed to update pet with ID %d: %v", id, err)
		http.Error(w, "failed to update pet", http.StatusInternalServerError)
		return
	}

	logger.InfoCtx(r.Context(), "Successfully updated pet with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Pet updated successfully",
	})
}

// DeletePet deletes a pet by ID
func DeletePet(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Deleting pet")
	
	// Get ID from URL
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid pet ID format for deletion: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Check if pet exists
	_, err = data.GetPetByID(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Pet not found for deletion with ID %d: %v", id, err)
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	// Delete pet
	err = data.DeletePet(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to delete pet with ID %d: %v", id, err)
		http.Error(w, "failed to delete pet", http.StatusInternalServerError)
		return
	}

	logger.InfoCtx(r.Context(), "Successfully deleted pet with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// -------------------- Pets --------------------

func GetPets(w http.ResponseWriter, r *http.Request) {
	logger.DebugCtx(r.Context(), "Fetching all pets")
	rows, err := data.ListPets(DB)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to fetch pets: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	pets := []Pet{}
	for _, rp := range rows {
		pets = append(pets, Pet{ID: rp.ID, Name: rp.Name, Species: rp.Species, Breed: rp.Breed, Birth: rp.Birth, OwnerID: rp.OwnerID})
	}
	logger.DebugCtx(r.Context(), "Retrieved %d pets", len(pets))
	json.NewEncoder(w).Encode(pets)
}

func GetPetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid pet ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Fetching pet with ID: %d", id)
	rp, err := data.GetPetByID(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Pet not found with ID %d: %v", id, err)
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	p := Pet{ID: rp.ID, Name: rp.Name, Species: rp.Species, Breed: rp.Breed, Birth: rp.Birth, OwnerID: rp.OwnerID}
	logger.DebugCtx(r.Context(), "Successfully retrieved pet: %+v", p)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func CreatePet(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Creating new pet")
	var p Pet
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		logger.WarnCtx(r.Context(), "Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Processing pet data: %+v", p)
	id, err := data.CreatePet(DB, data.PetInput{
		Name:    p.Name,
		Species: p.Species,
		Breed:   p.Breed,
		Birth:   p.Birth,
		OwnerID: p.OwnerID,
	})

	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to create pet: %v", err)
		http.Error(w, "failed to create pet", http.StatusInternalServerError)
		return
	}

	p.ID = id
	logger.InfoCtx(r.Context(), "Successfully created pet with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// -------------------- Vets --------------------

func GetVets(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Fetching all vets")
	rows, err := data.ListVets(DB)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to fetch vets: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	vets := make([]Vet, len(rows))
	for i, rv := range rows {
		vets[i] = Vet{
			ID:            rv.ID,
			Name:          rv.Name,
			Specialization: rv.Specialization,
		}
	}

	logger.DebugCtx(r.Context(), "Retrieved %d vets", len(vets))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vets)
}

func GetVetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid vet ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Fetching vet with ID: %d", id)
	v, err := data.GetVetByID(DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WarnCtx(r.Context(), "Vet not found with ID %d", id)
			http.Error(w, "vet not found", http.StatusNotFound)
		} else {
			logger.ErrorCtx(r.Context(), "Error fetching vet with ID %d: %v", id, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	vet := Vet{
		ID:            v.ID,
		Name:          v.Name,
		Specialization: v.Specialization,
	}

	logger.DebugCtx(r.Context(), "Successfully retrieved vet: %+v", vet)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vet)
}

func CreateVet(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Creating new vet")
	var v Vet
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logger.WarnCtx(r.Context(), "Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if v.Name == "" {
		logger.WarnCtx(r.Context(), "Vet name is required")
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Processing vet data: %+v", v)
	id, err := data.CreateVet(DB, data.VetInput{
		Name:          v.Name,
		Specialization: v.Specialization,
	})

	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to create vet: %v", err)
		http.Error(w, "failed to create vet", http.StatusInternalServerError)
		return
	}

	v.ID = id
	logger.InfoCtx(r.Context(), "Successfully created vet with ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func UpdateVet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid vet ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Check if vet exists
	_, err = data.GetVetByID(DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WarnCtx(r.Context(), "Vet not found with ID %d", id)
			http.Error(w, "vet not found", http.StatusNotFound)
		} else {
			logger.ErrorCtx(r.Context(), "Error fetching vet with ID %d: %v", id, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	var v Vet
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logger.WarnCtx(r.Context(), "Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if v.Name == "" {
		logger.WarnCtx(r.Context(), "Vet name is required")
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Updating vet ID %d with data: %+v", id, v)
	err = data.UpdateVet(DB, id, data.VetInput{
		Name:          v.Name,
		Specialization: v.Specialization,
	})

	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to update vet ID %d: %v", id, err)
		http.Error(w, "failed to update vet", http.StatusInternalServerError)
		return
	}

	// Fetch updated vet to return
	updatedVet, _ := data.GetVetByID(DB, id)
	vet := Vet{
		ID:            updatedVet.ID,
		Name:          updatedVet.Name,
		Specialization: updatedVet.Specialization,
	}

	logger.InfoCtx(r.Context(), "Successfully updated vet ID: %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vet)
}

func DeleteVet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid vet ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Check if vet exists
	_, err = data.GetVetByID(DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WarnCtx(r.Context(), "Vet not found with ID %d", id)
			http.Error(w, "vet not found", http.StatusNotFound)
		} else {
			logger.ErrorCtx(r.Context(), "Error fetching vet with ID %d: %v", id, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	logger.InfoCtx(r.Context(), "Deleting vet with ID: %d", id)
	err = data.DeleteVet(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to delete vet ID %d: %v", id, err)
		http.Error(w, "failed to delete vet", http.StatusInternalServerError)
		return
	}

	logger.InfoCtx(r.Context(), "Successfully deleted vet with ID: %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// -------------------- Visits --------------------

func GetVisits(w http.ResponseWriter, r *http.Request) {
	logger.DebugCtx(r.Context(), "Fetching all visits")
	rows, err := data.ListVisits(DB)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to fetch visits: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	visits := []Visit{}
	for _, rv := range rows {
		visits = append(visits, Visit{ID: rv.ID, PetID: rv.PetID, VetID: rv.VetID, Visit: rv.Visit, Desc: rv.Desc})
	}
	logger.DebugCtx(r.Context(), "Retrieved %d visits", len(visits))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(visits)
}

func GetVisitByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnCtx(r.Context(), "Invalid visit ID format: %s", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Fetching visit with ID: %d", id)
	rv, err := data.GetVisitByID(DB, id)
	if err != nil {
		logger.ErrorCtx(r.Context(), "Visit not found with ID %d: %v", id, err)
		http.Error(w, "visit not found", http.StatusNotFound)
		return
	}

	v := Visit{ID: rv.ID, PetID: rv.PetID, VetID: rv.VetID, Visit: rv.Visit, Desc: rv.Desc}
	logger.DebugCtx(r.Context(), "Successfully retrieved visit: %+v", v)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func CreateVisit(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context(), "Creating new visit")
	var v Visit
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logger.WarnCtx(r.Context(), "Invalid JSON in request: %v", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	logger.DebugCtx(r.Context(), "Processing visit data: %+v", v)
	id, err := data.CreateVisit(DB, data.VisitInput{
		PetID: v.PetID,
		VetID: v.VetID,
		Visit: v.Visit,
		Desc:  v.Desc,
	})

	if err != nil {
		logger.ErrorCtx(r.Context(), "Failed to create visit: %v", err)
		http.Error(w, "failed to create visit", http.StatusInternalServerError)
		return
	}

	v.ID = id
	logger.InfoCtx(r.Context(), "Successfully created visit with ID: %d", id)
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
