package data

import "database/sql"

type OwnerRow struct {
	ID      int
	Name    string
	Phone   string
	Address string
}

type OwnerInput struct {
	Name    string
	Phone   string
	Address string
}

func ListOwners(db *sql.DB) ([]OwnerRow, error) {
	rows, err := db.Query("SELECT id, name, phone, address FROM owners")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []OwnerRow{}
	for rows.Next() {
		var o OwnerRow
		if err := rows.Scan(&o.ID, &o.Name, &o.Phone, &o.Address); err != nil {
			return nil, err
		}
		res = append(res, o)
	}
	return res, nil
}

func GetOwnerByID(db *sql.DB, id int) (OwnerRow, error) {
	var o OwnerRow
	err := db.QueryRow("SELECT id, name, phone, address FROM owners WHERE id=$1", id).
		Scan(&o.ID, &o.Name, &o.Phone, &o.Address)
	return o, err
}

func CreateOwner(db *sql.DB, in OwnerInput) (int, error) {
	var id int
	err := db.QueryRow(
		"INSERT INTO owners(name,phone,address) VALUES($1,$2,$3) RETURNING id",
		in.Name, in.Phone, in.Address,
	).Scan(&id)
	return id, err
}

// UpdateOwner updates an existing owner in the database
func UpdateOwner(db *sql.DB, id int, in OwnerInput) error {
	sqlStatement := `
		UPDATE owners 
		SET name = $1, phone = $2, address = $3
		WHERE id = $4`
	
	_, err := db.Exec(sqlStatement, in.Name, in.Phone, in.Address, id)
	return err
}

// DeleteOwner removes an owner from the database
func DeleteOwner(db *sql.DB, id int) error {
	sqlStatement := `DELETE FROM owners WHERE id = $1`
	_, err := db.Exec(sqlStatement, id)
	return err
}
