package supercharger

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type LocationList []string

func (ll LocationList) Value() (driver.Value, error) {
	bytes, err := json.Marshal(ll)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (ll *LocationList) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}

	err := json.Unmarshal(asBytes, &ll)
	if err != nil {
		return errors.New("Scan could not unmarshal to []string")
	}

	return nil
}

type PhoneList []Phone

func (pl PhoneList) Value() (driver.Value, error) {
	bytes, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (pl *PhoneList) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}

	err := json.Unmarshal(asBytes, &pl)
	if err != nil {
		return errors.New("Scan could not unmarshal to []Phone")
	}

	return nil
}

type Phone struct {
	Label  string `json:"label"`
	Number string `json:"number"`
}

func (p Phone) Value() (driver.Value, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

type EmailList []Email

func (el EmailList) Value() (driver.Value, error) {
	bytes, err := json.Marshal(el)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func (el *EmailList) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}

	err := json.Unmarshal(asBytes, &el)
	if err != nil {
		return errors.New("Scan could not unmarshal to []Email")
	}

	return nil
}

type Email struct {
	Label string `json:"label"`
	Email string `json:"email"`
}

func (e Email) Value() (driver.Value, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

// JSONBool converts many raw type of boolean representations into their
// primitive equivilant.
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
