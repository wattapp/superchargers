package supercharger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"

	"github.com/dewski/spatial"
)

const chargersURL = "https://www.tesla.com/findus"

var ErrNoSuperchargersFound = errors.New("No superchargers found")

type Supercharger struct {
	Address                string         `db:"address" json:"address"` // not null
	AddressLine1           *string        `db:"address_line_1" json:"address_line_1"`
	AddressLine2           *string        `db:"address_line_2" json:"address_line_2"`
	AddressNotes           *string        `db:"address_notes" json:"address_notes,omitempty"`
	Amenities              *string        `db:"amentities" json:"amentities,omitempty"`
	BaiduLat               *float64       `db:"-" json:"baidu_lat,string"`
	BaiduLng               *float64       `db:"-" json:"baidu_lng,string"`
	BaiduGeo               *spatial.Point `db:"baidu_geo" json:"baidu_geo"`
	Chargers               *string        `db:"chargers" json:"chargers,omitempty"`
	City                   string         `db:"city" json:"city"` // not null
	CommonName             string         `db:"common_name" json:"common_name"`
	Country                string         `db:"country" json:"country"` // not null
	DestinationChargerLogo *string        `db:"destination_charger_logo" json:"destination_charger_logo,omitempty"`
	DestinationWebsite     *string        `db:"destination_website" json:"destination_website,omitempty"`
	DirectionsLink         *string        `db:"directions_link" json:"directions_link,omitempty"`
	Emails                 EmailList      `db:"emails" json:"emails,omitempty"`
	Geocode                string         `db:"geocode" json:"geocode"` // not null
	Hours                  *string        `db:"hours" json:"hours,omitempty"`
	IsGallery              JSONBool       `db:"is_gallery" json:"is_gallery"` // not null
	KioskPinX              *int64         `db:"kiosk_pin_x" json:"kiosk_pin_x,string,omitempty"`
	KioskPinY              *int64         `db:"kiosk_pin_y" json:"kiosk_pin_y,string,omitempty"`
	KioskZoomPinX          *int64         `db:"kiosk_zoom_pin_x" json:"kiosk_zoom_pin_x,string,omitempty"`
	KioskZoomPinY          *int64         `db:"kiosk_zoom_pin_y" json:"kiosk_zoom_pin_y,string,omitempty"`
	Geo                    spatial.Point  `db:"geo" json:"geo"`
	Latitude               float64        `db:"-" json:"latitude,string"`
	Longitude              float64        `db:"-" json:"longitude,string"`
	LocationID             string         `db:"location_id" json:"location_id"`     // not null
	LocationType           LocationList   `db:"location_type" json:"location_type"` // not null
	Nid                    int64          `db:"nid" json:"nid,string"`              // not null
	OpenSoon               JSONBool       `db:"open_soon" json:"open_soon"`         // not null
	Path                   string         `db:"path" json:"path"`                   // not null
	PostalCode             *string        `db:"postal_code" json:"postal_code,omitempty"`
	ProvinceState          *string        `db:"province_state" json:"province_state,omitempty"`
	Region                 string         `db:"region" json:"region,omitempty"` // not null
	SalesPhone             PhoneList      `db:"sales_phone" json:"sales_phone,omitempty"`
	SalesRepresentative    JSONBool       `db:"sales_representative" json:"sales_representative,omitempty"`
	SubRegion              *string        `db:"sub_region" json:"sub_region,omitempty"`
	Title                  string         `db:"title" json:"title"` // not null
}

func (sc *Supercharger) UnmarshalJSON(data []byte) error {
	type Alias Supercharger
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(sc),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	sc.Geo = spatial.Point{
		Lat: aux.Latitude,
		Lng: aux.Longitude,
	}

	if aux.BaiduLat != nil && aux.BaiduLng != nil {
		sc.BaiduGeo = &spatial.Point{
			Lat: *aux.BaiduLat,
			Lng: *aux.BaiduLng,
		}
	}

	return nil
}

func Superchargers() ([]Supercharger, error) {
	resp, err := http.Get(chargersURL)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Received bad status")
	}

	stripRe := regexp.MustCompile(`\r?\n`)
	body := stripRe.ReplaceAllString(string(b), " ")

	// Looking for the location data in the response
	locationRe := regexp.MustCompile(`var location_data =\s+?(?P<json>\[.*\])\;`)
	output := locationRe.FindStringSubmatch(body)

	if len(output) != 2 {
		return nil, ErrNoSuperchargersFound
	}

	var superchargers []Supercharger
	err = json.Unmarshal([]byte(output[1]), &superchargers)
	if err != nil {
		return nil, err
	}

	return superchargers, nil
}

func (s Supercharger) Equal(b Supercharger) bool {
	if s.Nid != b.Nid {
		panic("Checking equality for the wrong type")
	}

	if s.Address != b.Address {
		fmt.Printf("Address updated\nbefore=%v\nafter= %v\n", s.Address, b.Address)
		return false
	}
	if !stringPointersEqual(s.AddressLine1, b.AddressLine1) {
		fmt.Printf("AddressLine1 updated\nbefore=%v\nafter= %v\n", *s.AddressLine1, *b.AddressLine1)
		return false
	}
	if !stringPointersEqual(s.AddressLine2, b.AddressLine2) {
		fmt.Printf("AddressLine2 updated\nbefore=%v\nafter= %v\n", *s.AddressLine2, *b.AddressLine2)
		return false
	}
	if !stringPointersEqual(s.AddressNotes, b.AddressNotes) {
		fmt.Printf("AddressNotes updated\nbefore=%v\nafter= %v\n", *s.AddressNotes, *b.AddressNotes)
		return false
	}
	if !stringPointersEqual(s.Amenities, b.Amenities) {
		fmt.Printf("Amenities updated\nbefore=%v\nafter= %v\n", s.Amenities, b.Amenities)
		return false
	}
	if !spatialPointersEqual(s.BaiduGeo, b.BaiduGeo) {
		fmt.Printf("BaiduGeo updated\nbefore=%+v\nafter= %+v\n", s.BaiduGeo, b.BaiduGeo)
		return false
	}
	if !stringPointersEqual(s.Chargers, b.Chargers) {
		fmt.Printf("Chargers updated\nbefore=%v\nafter= %v\n", s.Chargers, b.Chargers)
		return false
	}
	if s.City != b.City {
		fmt.Printf("City updated\nbefore=%v\nafter= %v\n", s.City, b.City)
		return false
	}
	if s.CommonName != b.CommonName {
		fmt.Printf("CommonName updated\nbefore=%v\nafter= %v\n", s.CommonName, b.CommonName)
		return false
	}
	if s.Country != b.Country {
		fmt.Printf("Country updated\nbefore=%v\nafter= %v\n", s.Country, b.Country)
		return false
	}
	if !stringPointersEqual(s.DestinationChargerLogo, b.DestinationChargerLogo) {
		fmt.Printf("DestinationChargerLogo updated\nbefore=%v\nafter= %v\n", s.DestinationChargerLogo, b.DestinationChargerLogo)
		return false
	}
	if !stringPointersEqual(s.DestinationWebsite, b.DestinationWebsite) {
		fmt.Printf("DestinationWebsite updated\nbefore=%v\nafter= %v\n", s.DestinationWebsite, b.DestinationWebsite)
		return false
	}
	if !stringPointersEqual(s.DirectionsLink, b.DirectionsLink) {
		fmt.Printf("DirectionsLink updated\nbefore=%v\nafter= %v\n", s.DirectionsLink, *b.DirectionsLink)
		return false
	}
	if s.Geocode != b.Geocode {
		fmt.Printf("Geocode updated\nbefore=%v\nafter= %v\n", s.Geocode, b.Geocode)
		return false
	}
	if !stringPointersEqual(s.Hours, b.Hours) {
		fmt.Printf("Hours updated\nbefore=%v\nafter= %v\n", s.Hours, b.Hours)
		return false
	}
	if s.IsGallery != b.IsGallery {
		fmt.Printf("IsGallery updated\nbefore=%v\nafter= %v\n", s.IsGallery, b.IsGallery)
		return false
	}
	if !int64PointersEqual(s.KioskPinX, b.KioskPinX) {
		fmt.Printf("KioskPinX updated\nbefore=%v\nafter= %v\n", s.KioskPinX, b.KioskPinX)
		return false
	}
	if !int64PointersEqual(s.KioskPinY, b.KioskPinY) {
		fmt.Printf("KioskPinY updated\nbefore=%v\nafter= %v\n", s.KioskPinY, b.KioskPinY)
		return false
	}
	if !int64PointersEqual(s.KioskZoomPinX, b.KioskZoomPinX) {
		fmt.Printf("KioskZoomPinX updated\nbefore=%v\nafter= %v\n", s.KioskZoomPinX, b.KioskZoomPinX)
		return false
	}
	if !int64PointersEqual(s.KioskZoomPinY, b.KioskZoomPinY) {
		fmt.Printf("KioskZoomPinY updated\nbefore=%v\nafter= %v\n", s.KioskZoomPinY, b.KioskZoomPinY)
		return false
	}
	if s.Geo != b.Geo {
		fmt.Printf("Geo updated\nbefore=%v\nafter= %v\n", s.Geo, b.Geo)
		return false
	}
	if s.LocationID != b.LocationID {
		fmt.Printf("LocationID updated\nbefore=%v\nafter= %v\n", s.LocationID, b.LocationID)
		return false
	}
	if s.OpenSoon != b.OpenSoon {
		fmt.Printf("OpenSoon updated\nbefore=%v\nafter= %v\n", s.OpenSoon, b.OpenSoon)
		return false
	}
	if s.Path != b.Path {
		fmt.Printf("Path updated\nbefore=%v\nafter= %v\n", s.Path, b.Path)
		return false
	}
	if !stringPointersEqual(s.PostalCode, b.PostalCode) {
		fmt.Printf("PostalCode updated\nbefore=%v\nafter= %v\n", s.PostalCode, b.PostalCode)
		return false
	}
	if !stringPointersEqual(s.ProvinceState, b.ProvinceState) {
		fmt.Printf("ProvinceState updated\nbefore=%v\nafter= %v\n", s.ProvinceState, b.ProvinceState)
		return false
	}
	if s.Region != b.Region {
		fmt.Printf("Region updated\nbefore=%v\nafter= %v\n", s.Region, b.Region)
		return false
	}
	if s.SalesRepresentative != b.SalesRepresentative {
		fmt.Printf("SalesRepresentative updated\nbefore=%v\nafter= %v\n", s.SalesRepresentative, b.SalesRepresentative)
		return false
	}
	if !stringPointersEqual(s.SubRegion, b.SubRegion) {
		fmt.Printf("SubRegion updated\nbefore=%v\nafter= %v\n", s.SubRegion, b.SubRegion)
		return false
	}
	if s.Title != b.Title {
		fmt.Printf("Title updated\nbefore=%v\nafter= %v\n", s.Title, b.Title)
		return false
	}

	if !reflect.DeepEqual(s.Emails, b.Emails) {
		fmt.Printf("Emails updated\nbefore=%v\nafter= %v\n", s.Emails, b.Emails)
		return false
	}

	if !reflect.DeepEqual(s.LocationType, b.LocationType) {
		fmt.Printf("LocationType updated\nbefore=%v\nafter= %v\n", s.LocationType, b.LocationType)
		return false
	}

	if !reflect.DeepEqual(s.SalesPhone, b.SalesPhone) {
		fmt.Printf("SalesPhone updated\nbefore=%v\nafter= %v\n", s.SalesPhone, b.SalesPhone)
		return false
	}

	return true
}

func stringPointersEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	return *a == *b
}

func int64PointersEqual(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	return *a == *b
}

func spatialPointersEqual(a, b *spatial.Point) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	return *a == *b
}
