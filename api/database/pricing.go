package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Pricing struct {
	Monday    []float64 `json:"monday"`
	Tuesday   []float64 `json:"tuesday"`
	Wednesday []float64 `json:"wednesday"`
	Thursday  []float64 `json:"thursday"`
	Friday    []float64 `json:"friday"`
	Saturday  []float64 `json:"saturday"`
	Sunday    []float64 `json:"sunday"`
}

// Value Marshal
func (a Pricing) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Pricing) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
