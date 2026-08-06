package dto

type CoordinatesInput struct {
	Lat float32 `json:"lat" minimum:"-90" maximum:"90" doc:"Latitude"`
	Lng float32 `json:"lng" minimum:"-180" maximum:"180" doc:"Longitude"`
}
