package location

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/graphql-go/relay"
	"github.com/wattapp/superchargers/pkg/database"
	"github.com/wattapp/superchargers/pkg/supercharger"
)

type Location struct {
	supercharger.Supercharger

	ID        int64     `db:"id" json:"id"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (l Location) Cursor() relay.ConnectionCursor {
	str := fmt.Sprintf("%v%v", relay.PREFIX, l.ID)
	return relay.ConnectionCursor(base64.StdEncoding.EncodeToString([]byte(str)))
}

func (l Location) ToGlobalID() string {
	id := strconv.FormatInt(l.ID, 10)
	return relay.ToGlobalID("Location", id)
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

func Locations() ([]*Location, error) {
	locations := []*Location{}
	err := database.Conn().
		Select("*").
		From("locations").
		QueryStructs(&locations)

	if err != nil {
		return nil, err
	}

	return locations, nil
}

func Update() {
	for {
		fmt.Println("Checking for super chargers...")
		err := Sync()
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(24 * time.Hour)
	}
}

func Sync() error {
	locations, err := supercharger.Superchargers()
	if err != nil {
		return err
	}

	for _, location := range locations {
		_, err = syncLocation(location)
		if err != nil {
			return err
		}
	}

	return nil
}

func syncLocation(sc supercharger.Supercharger) (*Location, error) {
	fmt.Printf("Looking up record for nid=%d\n", sc.Nid)
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
		fmt.Printf("Found record %d\n", location.ID)
		return location, nil
	}

	fmt.Printf("No record found for %d, preparing to create one\n", sc.Nid)

	location = &Location{Supercharger: sc}
	location.CreatedAt = time.Now().UTC()
	location.UpdatedAt = time.Now().UTC()
	err = database.Conn().
		InsertInto("locations").
		Columns(
			"address",
			"address_line_1",
			"address_line_2",
			"address_notes",
			"amentities",
			"baidu_lat",
			"baidu_lng",
			"chargers",
			"city",
			"common_name",
			"country",
			"destination_charger_logo",
			"destination_website",
			"emails",
			"geocode",
			"hours",
			"is_gallery",
			"kiosk_pin_x",
			"kiosk_pin_y",
			"kiosk_zoom_pin_x",
			"kiosk_zoom_pin_y",
			"latitude",
			"longitude",
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
		).
		Record(sc).
		Returning("id").
		QueryScalar(&location.ID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Created location %d for remote object %d at %v\n", location.ID, location.Nid, location.CreatedAt)

	return location, nil
}
