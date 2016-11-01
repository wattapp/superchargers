
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX index_locations_on_region ON locations(region);
CREATE INDEX index_locations_on_country ON locations(country);
CREATE INDEX index_locations_on_location_type ON locations USING GIN(location_type);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX index_locations_on_location_type;
DROP INDEX index_locations_on_region;
DROP INDEX index_locations_on_country;
