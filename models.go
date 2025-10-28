//go:build ignore
// +build ignore

package main

import "time"

type Owner struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type Pet struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Species  string    `json:"species"`
	Breed    string    `json:"breed"`
	Birth    time.Time `json:"birth_date"`
	OwnerID  int       `json:"owner_id"`
}

type Visit struct {
	ID     int       `json:"id"`
	PetID  int       `json:"pet_id"`
	VetID  int       `json:"vet_id"`
	Visit  time.Time `json:"visit_date"`
	Desc   string    `json:"description"`
}
