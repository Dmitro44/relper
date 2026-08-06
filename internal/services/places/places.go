package places

import (
	"context"

	"github.com/Dmitro44/relper/internal/entities/coordinates"
	"github.com/Dmitro44/relper/internal/entities/place"
	placesrepo "github.com/Dmitro44/relper/internal/repositories/places"
	"github.com/Dmitro44/relper/internal/services/locations"
	"github.com/google/uuid"
)

type PlaceService struct {
	repo            placesrepo.Repository
	locationservice locations.LocationService
}

func NewPlaceService(repo placesrepo.Repository, locationservice locations.LocationService) *PlaceService {
	return &PlaceService{
		repo:            repo,
		locationservice: locationservice,
	}
}

type CreateInput struct {
	LocationID  uuid.UUID
	Name        string
	Category    string
	Address     string
	Phone       string
	WebsiteURL  string
	Coordinates coordinates.Coordinates
	Description string
	Rating      float32
	Upvotes     int32
}

type UpdateInput struct {
	LocationID  *uuid.UUID
	Name        *string
	Category    *string
	Address     *string
	Phone       *string
	WebsiteURL  *string
	Coordinates *coordinates.Coordinates
	Description *string
	Rating      *float32
	Upvotes     *int32
}

func (p *PlaceService) GetByID(ctx context.Context, id uuid.UUID) (*place.Place, error) {
	placeEntity, err := p.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return placeEntity, nil
}

func (p *PlaceService) Create(ctx context.Context, input CreateInput) (*place.Place, error) {
	// check if location id exists
	_, err := p.locationservice.GetByID(ctx, input.LocationID)
	if err != nil {
		return nil, err
	}

	placeEntity, err := p.repo.Create(ctx, placesrepo.CreateInput{
		LocationID:  input.LocationID,
		Name:        input.Name,
		Category:    input.Category,
		Address:     input.Address,
		Phone:       input.Phone,
		WebsiteURL:  input.WebsiteURL,
		Coordinates: input.Coordinates,
		Description: input.Description,
		Rating:      input.Rating,
		Upvotes:     input.Upvotes,
	})
	if err != nil {
		return nil, err
	}

	return placeEntity, nil
}

func (p *PlaceService) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*place.Place, error) {
	placeEntity, err := p.repo.Update(ctx, id, placesrepo.UpdateInput{
		LocationID:  input.LocationID,
		Name:        input.Name,
		Category:    input.Category,
		Address:     input.Address,
		Phone:       input.Phone,
		WebsiteURL:  input.WebsiteURL,
		Coordinates: input.Coordinates,
		Description: input.Description,
		Rating:      input.Rating,
		Upvotes:     input.Upvotes,
	})
	if err != nil {
		return nil, err
	}

	return placeEntity, nil
}

func (p *PlaceService) Delete(ctx context.Context, id uuid.UUID) error {
	err := p.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
