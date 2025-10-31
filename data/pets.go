package data

import (
	"database/sql"
	"time"
)

type PetRow struct {
	ID      int
	Name    string
	Species string
	Breed   string
	Birth   time.Time
	OwnerID int
}

type PetInput struct {
	Name    string
	Species string
	Breed   string
	Birth   time.Time
	OwnerID int
}

func ListPets(db *sql.DB) ([]PetRow, error) {
	rows, err := db.Query("SELECT id, name, species, breed, birth_date, owner_id FROM pets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	s := []PetRow{}
	for rows.Next() {
		var p PetRow
		if err := rows.Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.Birth, &p.OwnerID); err != nil {
			return nil, err
		}
		s = append(s, p)
	}
	return s, nil
}

func GetPetByID(db *sql.DB, id int) (PetRow, error) {
	var p PetRow
	err := db.QueryRow("SELECT id, name, species, breed, birth_date, owner_id FROM pets WHERE id=$1", id).
		Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.Birth, &p.OwnerID)
	return p, err
}

func CreatePet(db *sql.DB, in PetInput) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO pets(name,species,breed,birth_date,owner_id) VALUES($1,$2,$3,$4,$5) RETURNING id",
		in.Name, in.Species, in.Breed, in.Birth, in.OwnerID,
	).Scan(&id)
	return id, err
}
