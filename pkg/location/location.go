package location

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dewski/spatial"
	"github.com/graphql-go/relay"
	"github.com/wattapp/superchargers/pkg/database"
	"github.com/wattapp/superchargers/pkg/supercharger"
)

var columns = []string{
	"address",
	"address_line_1",
	"address_line_2",
	"address_notes",
	"amentities",
	"chargers",
	"city",
	"common_name",
	"country",
	"destination_charger_logo",
	"destination_website",
	"directions_link",
	"emails",
	"geocode",
	"hours",
	"is_gallery",
	"kiosk_pin_x",
	"kiosk_pin_y",
	"kiosk_zoom_pin_x",
	"kiosk_zoom_pin_y",
	"geo",
	"location_id",
	"location_type",
	"nid",
	"open_soon",
	"path",
	"postal_code",
	"province_state",
	"region",
	"sales_phone",
	"sales_representative",
	"sub_region",
	"title",
	"updated_at",
	"created_at",
}

type Location struct {
	supercharger.Supercharger

	ID        int64     `db:"id" json:"id"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func GetLocation(locationID int64) (*Location, error) {
	location := &Location{}
	err := database.Conn().
		Select("*").
		From("locations").
		Where("id = $1", locationID).
		QueryStruct(location)

	if err != nil {
		return nil, err
	}

	return location, nil
}

func Near(scope database.GraphQLScope) ([]*Location, error) {
	lat, ok := scope.Args["latitude"].(float64)
	if !ok {
		return nil, errors.New("Invalid latitude")
	}

	lng, ok := scope.Args["longitude"].(float64)
	if !ok {
		return nil, errors.New("Invalid longitude")
	}

	point := spatial.Point{
		Lat: lat,
		Lng: lng,
	}

	locations := []*Location{}
	builder := database.Conn().
		Select("*").
		From("locations").
		OrderBy("geo <-> $1::geometry", point)

	if scope.ConnectionArguments.First != -1 {
		builder = builder.Limit(uint64(scope.ConnectionArguments.First))
	} else {
		builder = builder.Limit(uint64(database.DefaultLimit))
	}

	err := builder.QueryStructs(&locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func Locations(scope database.GraphQLScope) ([]*Location, error) {
	locations := []*Location{}
	builder := database.Conn().
		Select("*").
		From("locations")

	if scope.Args["region"] != nil {
		var regions []string
		for _, r := range scope.Args["region"].([]interface{}) {
			regions = append(regions, r.(string))
		}

		if len(regions) > 0 {
			builder = builder.Where("region IN $1", regions)
		}
	}

	if scope.Args["country"] != nil {
		var countries []string
		for _, c := range scope.Args["country"].([]interface{}) {
			countries = append(countries, c.(string))
		}

		if len(countries) > 0 {
			builder = builder.Where("country IN $1", countries)
		}
	}

	if scope.Args["openSoon"] != nil {
		builder = builder.Where("open_soon = $1", scope.Args["openSoon"])
	}

	if scope.Args["isGallery"] != nil {
		builder = builder.Where("is_gallery = $1", scope.Args["isGallery"])
	}

	bb, ok := scope.Args["boundingBox"].([]interface{})
	if ok && len(bb) == 4 {
		nwLat, nwLng, seLat, seLng := bb[0], bb[1], bb[2], bb[3]
		builder = builder.Where(
			`ST_Contains(ST_SetSRID(ST_MakeBox2D(ST_Point($1, $2), ST_Point($3, $4)), 4326), geo)`,
			nwLng,
			nwLat,
			seLng,
			seLat,
		)
	}

	if scope.Args["type"] != nil {
		var types []string
		for _, t := range scope.Args["type"].([]interface{}) {
			types = append(types, t.(string))
		}

		if len(types) > 0 {
			builder = builder.
				SetIsInterpolated(false).
				Where(fmt.Sprintf("location_type ?| array['%s']", strings.Join(types, "','"))).
				SetIsInterpolated(true)
		}
	}

	scope.OrderBy = database.OrderOnCreatedAt
	query, err := database.ApplyGraphQLScope(builder, scope)
	if err != nil {
		return nil, err
	}

	err = query.QueryStructs(&locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (l Location) Cursor() relay.ConnectionCursor {
	str := fmt.Sprintf("%v%v", relay.PREFIX, l.ID)
	return relay.ConnectionCursor(base64.StdEncoding.EncodeToString([]byte(str)))
}

func (l Location) ToGlobalID() string {
	id := strconv.FormatInt(l.ID, 10)
	return relay.ToGlobalID("Location", id)
}

func (l Location) Update(sc supercharger.Supercharger) error {
	l.UpdatedAt = time.Now().UTC()
	l.Supercharger = sc
	_, err := database.Conn().
		Update("locations").
		Set("address", l.Address).
		Set("address_line_1", l.AddressLine1).
		Set("address_line_2", l.AddressLine2).
		Set("address_notes", l.AddressNotes).
		Set("amentities", l.Amenities).
		Set("chargers", l.Chargers).
		Set("city", l.City).
		Set("common_name", l.CommonName).
		Set("country", l.Country).
		Set("destination_charger_logo", l.DestinationChargerLogo).
		Set("destination_website", l.DestinationWebsite).
		Set("directions_link", l.DirectionsLink).
		Set("emails", l.Emails).
		Set("geocode", l.Geocode).
		Set("hours", l.Hours).
		Set("is_gallery", l.IsGallery).
		Set("kiosk_pin_x", l.KioskPinX).
		Set("kiosk_pin_y", l.KioskPinY).
		Set("kiosk_zoom_pin_x", l.KioskZoomPinX).
		Set("kiosk_zoom_pin_y", l.KioskZoomPinY).
		Set("geo", l.Geo).
		Set("location_id", l.LocationID).
		Set("location_type", l.LocationType).
		Set("nid", l.Nid).
		Set("open_soon", l.OpenSoon).
		Set("path", l.Path).
		Set("postal_code", l.PostalCode).
		Set("province_state", l.ProvinceState).
		Set("region", l.Region).
		Set("sales_phone", l.SalesPhone).
		Set("sales_representative", l.SalesRepresentative).
		Set("sub_region", l.SubRegion).
		Set("title", l.Title).
		Set("updated_at", l.UpdatedAt).
		Where("id = $1", l.ID).
		Exec()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully updated nid=%d\n", l.Nid)

	return nil
}

func Sync() (added, updated int, err error) {
	added, updated = 0, 0
	start := time.Now().UTC()
	locations, err := supercharger.Superchargers()
	if err != nil {
		return
	}

	for _, location := range locations {
		var l *Location
		l, err = syncLocation(location)
		if err != nil {
			return
		}

		if l.CreatedAt.After(start) {
			added += 1
		} else if l.UpdatedAt.After(start) {
			updated += 1
		}
	}

	return
}

func syncLocation(sc supercharger.Supercharger) (*Location, error) {
	location := &Location{}
	err := database.Conn().
		Select("*").
		From("locations").
		Where("nid = $1", sc.Nid).
		QueryStruct(location)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if location.ID > 0 {
		if !location.Supercharger.Equal(sc) {
			fmt.Printf("Remote record for nid=%d has been updated, updating in database\n", location.Nid)
			err = location.Update(sc)
			if err != nil {
				return nil, err
			}
		}

		return location, nil
	}

	fmt.Printf("No record found for %d, preparing to create one\n", sc.Nid)

	location = &Location{Supercharger: sc}
	if sc.BaiduLat != nil && sc.BaiduLng != nil && sc.Latitude == 0.0 && sc.Longitude == 0.0 {
		location.Geo = spatial.Point{
			Lat: *sc.BaiduLat,
			Lng: *sc.BaiduLng,
		}
	} else {
		location.Geo = spatial.Point{
			Lat: sc.Latitude,
			Lng: sc.Longitude,
		}
	}
	location.CreatedAt = time.Now().UTC()
	location.UpdatedAt = time.Now().UTC()
	err = database.Conn().
		InsertInto("locations").
		Columns(columns...).
		Record(location).
		Returning("id").
		QueryScalar(&location.ID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Created location %d for remote object %d at %v\n", location.ID, location.Nid, location.CreatedAt)

	return location, nil
}
