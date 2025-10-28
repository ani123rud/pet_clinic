package data

import (
	"petclinic/models"
)

func ListPets() ([]models.Pet, error) {
	rows, err := DB.Query("SELECT id, name, species, breed, birth_date, owner_id FROM pets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pets := []models.Pet{}
	for rows.Next() {
		var p models.Pet
		if err := rows.Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.Birth, &p.OwnerID); err != nil {
			return nil, err
		}
		pets = append(pets, p)
	}
	return pets, nil
}

func GetPetByID(id int) (models.Pet, error) {
	var p models.Pet
	err := DB.QueryRow("SELECT id, name, species, breed, birth_date, owner_id FROM pets WHERE id=$1", id).
		Scan(&p.ID, &p.Name, &p.Species, &p.Breed, &p.Birth, &p.OwnerID)
	return p, err
}

func CreatePet(p models.Pet) (models.Pet, error) {
	err := DB.QueryRow(
		"INSERT INTO pets(name,species,breed,birth_date,owner_id) VALUES($1,$2,$3,$4,$5) RETURNING id",
		p.Name, p.Species, p.Breed, p.Birth, p.OwnerID,
	).Scan(&p.ID)
	return p, err
}
