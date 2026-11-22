package dto

type CoordinatesInput struct {
	Latitude  float64 `json:"lat" minimum:"-90" maximum:"90" doc:"Latitude"`
	Longitude float64 `json:"lng" minimum:"-180" maximum:"180" doc:"Longitude"`
}
