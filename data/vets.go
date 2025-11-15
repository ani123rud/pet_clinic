package data

import "database/sql"

type VetRow struct {
	ID            int
	Name          string
	Specialization string
}

type VetInput struct {
	Name          string
	Specialization string
}

func ListVets(db *sql.DB) ([]VetRow, error) {
	rows, err := db.Query("SELECT id, name, specialization FROM vets ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vets []VetRow
	for rows.Next() {
		var v VetRow
		if err := rows.Scan(&v.ID, &v.Name, &v.Specialization); err != nil {
			return nil, err
		}
		vets = append(vets, v)
	}
	return vets, nil
}

func GetVetByID(db *sql.DB, id int) (VetRow, error) {
	var v VetRow
	err := db.QueryRow("SELECT id, name, specialization FROM vets WHERE id = $1", id).
		Scan(&v.ID, &v.Name, &v.Specialization)
	return v, err
}

func CreateVet(db *sql.DB, in VetInput) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO vets(name, specialization) VALUES($1, $2) RETURNING id",
		in.Name, in.Specialization,
	).Scan(&id)

	return id, err
}

func UpdateVet(db *sql.DB, id int, in VetInput) error {
	_, err := db.Exec(
		"UPDATE vets SET name = $1, specialization = $2 WHERE id = $3",
		in.Name, in.Specialization, id,
	)
	return err
}

func DeleteVet(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM vets WHERE id = $1", id)
	return err
}
