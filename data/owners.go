package data

import (
	"petclinic/models"
)

func ListOwners() ([]models.Owner, error) {
	rows, err := DB.Query("SELECT id, name, phone, address FROM owners")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	owners := []models.Owner{}
	for rows.Next() {
		var o models.Owner
		if err := rows.Scan(&o.ID, &o.Name, &o.Phone, &o.Address); err != nil {
			return nil, err
		}
		owners = append(owners, o)
	}
	return owners, nil
}

func GetOwnerByID(id int) (models.Owner, error) {
	var o models.Owner
	err := DB.QueryRow("SELECT id, name, phone, address FROM owners WHERE id=$1", id).
		Scan(&o.ID, &o.Name, &o.Phone, &o.Address)
	return o, err
}

func CreateOwner(o models.Owner) (models.Owner, error) {
	err := DB.QueryRow(
		"INSERT INTO owners(name,phone,address) VALUES($1,$2,$3) RETURNING id",
		o.Name, o.Phone, o.Address,
	).Scan(&o.ID)
	return o, err
}
