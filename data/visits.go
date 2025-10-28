package data

import (
	"petclinic/models"
)

func ListVisits() ([]models.Visit, error) {
	rows, err := DB.Query("SELECT id, pet_id, vet_id, visit_date, description FROM visits")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	visits := []models.Visit{}
	for rows.Next() {
		var v models.Visit
		if err := rows.Scan(&v.ID, &v.PetID, &v.VetID, &v.Visit, &v.Desc); err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	return visits, nil
}

func GetVisitByID(id int) (models.Visit, error) {
	var v models.Visit
	err := DB.QueryRow("SELECT id, pet_id, vet_id, visit_date, description FROM visits WHERE id=$1", id).
		Scan(&v.ID, &v.PetID, &v.VetID, &v.Visit, &v.Desc)
	return v, err
}

func CreateVisit(v models.Visit) (models.Visit, error) {
	err := DB.QueryRow(
		"INSERT INTO visits(pet_id,vet_id,visit_date,description) VALUES($1,$2,$3,$4) RETURNING id",
		v.PetID, v.VetID, v.Visit, v.Desc,
	).Scan(&v.ID)
	return v, err
}
