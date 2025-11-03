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

// UpdateVisit updates an existing visit in the database
func UpdateVisit(db *sql.DB, id int, in VisitInput) error {
	sqlStatement := `
		UPDATE visits 
		SET pet_id = $1, vet_id = $2, visit_date = $3, description = $4
		WHERE id = $5`
	
	_, err := db.Exec(sqlStatement, in.PetID, in.VetID, in.Visit, in.Desc, id)
	return err
}

// DeleteVisit removes a visit from the database
func DeleteVisit(db *sql.DB, id int) error {
	sqlStatement := `DELETE FROM visits WHERE id = $1`
	_, err := db.Exec(sqlStatement, id)
	return err
}
