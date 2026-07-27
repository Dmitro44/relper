package model

import "github.com/google/uuid"

type Location struct {
	ID          uuid.UUID
	Country     string
	City        string
	Coordinates Coordinates
}

type Coordinates struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}
