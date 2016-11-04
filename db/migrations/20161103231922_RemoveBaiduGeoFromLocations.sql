
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP INDEX index_locations_on_baidu_geo;
ALTER TABLE locations DROP COLUMN baidu_geo;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE locations ADD COLUMN baidu_geo geometry(Geometry, 4326) null;
CREATE INDEX index_locations_on_baidu_geo ON locations USING GIST(baidu_geo);
