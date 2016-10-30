package supercharger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const chargersURL = "https://www.tesla.com/findus"

var ErrNoSuperchargersFound = errors.New("No superchargers found")

type Phone struct {
	Label  string
	Number string
}

type Email struct {
	Label string
	Email string
}

type Location struct {
	Address                string   `json:"address"` // not null
	AddressLine1           string   `json:"address_line_1"`
	AddressLine2           string   `json:"address_line_2"`
	AddressNotes           string   `json:"address_notes,omitempty"`
	Amenities              string   `json:"amentities,omitempty"`
	BaiduLat               float64  `json:"baidu_lat,string"`
	BaiduLng               float64  `json:"baidu_lng,string"`
	Chargers               string   `json:"chargers,omitempty"`
	City                   string   `json:"city"` // not null
	CommonName             string   `json:"common_name"`
	Country                string   `json:"country"` // not null
	DestinationChargerLogo string   `json:"destination_charger_logo,omitempty"`
	DestinationWebsite     string   `json:"destination_website,omitempty"`
	DirectionsLink         string   `json:"directions_link,omitempty"`
	Emails                 []Email  `json:"emails,omitempty"`
	Geocode                string   `json:"geocode"` // not null
	Hours                  string   `json:"hours,omitempty"`
	IsGallery              JSONBool `json:"is_gallery"` // not null
	KioskPinX              int64    `json:"kiosk_pin_x,string,omitempty"`
	KioskPinY              int64    `json:"kiosk_pin_y,string,omitempty"`
	KioskZoomPinX          int64    `json:"kiosk_zoom_pin_x,string,omitempty"`
	KioskZoomPinY          int64    `json:"kiosk_zoom_pin_y,string,omitempty"`
	Latitude               float64  `json:"latitude,string"`
	Longitude              float64  `json:"longitude,string"`
	LocationID             string   `json:"location_id"`   // not null
	LocationType           []string `json:"location_type"` // not null
	Nid                    int64    `json:"nid,string"`    // not null
	OpenSoon               JSONBool `json:"open_soon"`     // not null
	Path                   string   `json:"path"`          // not null
	PostalCode             string   `json:"postal_code,omitempty"`
	ProvinceState          string   `json:"province_state,omitempty"`
	Region                 string   `json:"region,omitempty"` // not null
	SalesPhone             []Phone  `json:"sales_phone,omitempty"`
	SalesRepresentative    JSONBool `json:"sales_representative,omitempty"`
	SubRegion              string   `json:"sub_region,omitempty"`
	Title                  string   `json:"title"` // not null
}

type JSONBool bool

func (b *JSONBool) UnmarshalJSON(data []byte) error {
	s := string(data)

	switch s {
	case "1", "\"1\"", "true":
		*b = true
	case "0", "\"0\"", "false":
		*b = false
	default:
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", s))
	}

	return nil
}

func Superchargers() (locations []Location, err error) {
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

	err = json.Unmarshal([]byte(output[1]), &locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}
