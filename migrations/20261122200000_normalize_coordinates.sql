-- +goose Up
-- +goose StatementBegin
ALTER TABLE places
    ADD COLUMN latitude  double precision,
    ADD COLUMN longitude double precision;

UPDATE places
SET latitude  = (coordinates->>'lat')::double precision,
    longitude = (coordinates->>'lng')::double precision
WHERE coordinates IS NOT NULL;

ALTER TABLE places
    ALTER COLUMN latitude  SET NOT NULL,
    ALTER COLUMN longitude SET NOT NULL,
    ADD CONSTRAINT places_latitude_range  CHECK (latitude  BETWEEN -90  AND 90),
    ADD CONSTRAINT places_longitude_range CHECK (longitude BETWEEN -180 AND 180);

ALTER TABLE places DROP COLUMN coordinates;

ALTER TABLE locations
    ADD COLUMN latitude  double precision,
    ADD COLUMN longitude double precision;

UPDATE locations
SET latitude  = (coordinates->>'lat')::double precision,
    longitude = (coordinates->>'lng')::double precision
WHERE coordinates IS NOT NULL;

ALTER TABLE locations
    ALTER COLUMN latitude  SET NOT NULL,
    ALTER COLUMN longitude SET NOT NULL,
    ADD CONSTRAINT locations_latitude_range  CHECK (latitude  BETWEEN -90  AND 90),
    ADD CONSTRAINT locations_longitude_range CHECK (longitude BETWEEN -180 AND 180);

ALTER TABLE locations DROP COLUMN coordinates;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE locations
    ADD COLUMN coordinates jsonb;

UPDATE locations
SET coordinates = jsonb_build_object('lat', latitude, 'lng', longitude);

ALTER TABLE locations
    DROP CONSTRAINT locations_latitude_range,
    DROP CONSTRAINT locations_longitude_range,
    DROP COLUMN latitude,
    DROP COLUMN longitude;

ALTER TABLE places
    ADD COLUMN coordinates jsonb;

UPDATE places
SET coordinates = jsonb_build_object('lat', latitude, 'lng', longitude);

ALTER TABLE places
    DROP CONSTRAINT places_latitude_range,
    DROP CONSTRAINT places_longitude_range,
    DROP COLUMN latitude,
    DROP COLUMN longitude;
-- +goose StatementEnd
