
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION postgis;
CREATE TABLE locations (
  id serial primary key,
  address text not null,
  address_line_1 text null,
  address_line_2 text null,
  address_notes text null,
  amentities text null,
  baidu_geo geometry(Geometry, 4326) null,
  chargers text null,
  city varchar(100) not null,
  common_name text null,
  country varchar(50) not null,
  destination_charger_logo text null,
  destination_website text null,
  directions_link text null,
  emails jsonb null, -- array
  geocode varchar(100) not null,
  hours text null,
  is_gallery bool not null default false,
  kiosk_pin_x integer null,
  kiosk_pin_y integer null,
  kiosk_zoom_pin_x integer null,
  kiosk_zoom_pin_y integer null,
  geo geometry(Geometry, 4326) not null,
  location_id varchar(255) not null,
  location_type jsonb not null, -- array
  nid integer not null unique,
  open_soon bool not null default false,
  path varchar(100) not null,
  postal_code varchar(50) null,
  province_state varchar(100) null,
  region varchar(100) not null,
  sales_phone jsonb null, --array
  sales_representative bool not null default false,
  sub_region varchar(100) null,
  title varchar(255) not null,
  updated_at timestamp(3) not null,
  created_at timestamp(3) not null
);

CREATE INDEX index_locations_on_geo ON locations USING GIST(geo);
CREATE INDEX index_locations_on_baidu_geo ON locations USING GIST(baidu_geo);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE locations;
DROP EXTENSION postgis;
