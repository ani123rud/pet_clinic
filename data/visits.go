package data

import (
	"database/sql"
	"time"
)

type VisitRow struct {
	ID    int
	PetID int
	VetID int
	Visit time.Time
	Desc  string
}

type VisitInput struct {
	PetID int
	VetID int
	Visit time.Time
	Desc  string
}

func ListVisits(db *sql.DB) ([]VisitRow, error) {
	rows, err := db.Query("SELECT id, pet_id, vet_id, visit_date, description FROM visits")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []VisitRow{}
	for rows.Next() {
		var v VisitRow
		if err := rows.Scan(&v.ID, &v.PetID, &v.VetID, &v.Visit, &v.Desc); err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func GetVisitByID(db *sql.DB, id int) (VisitRow, error) {
	var v VisitRow
	err := db.QueryRow("SELECT id, pet_id, vet_id, visit_date, description FROM visits WHERE id=$1", id).
		Scan(&v.ID, &v.PetID, &v.VetID, &v.Visit, &v.Desc)
	return v, err
}

func CreateVisit(db *sql.DB, in VisitInput) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO visits(pet_id,vet_id,visit_date,description) VALUES($1,$2,$3,$4) RETURNING id",
		in.PetID, in.VetID, in.Visit, in.Desc,
	).Scan(&id)
	return id, err
}
