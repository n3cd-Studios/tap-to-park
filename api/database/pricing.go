package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
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

func (spot *Spot) GetPrice() float64 {
	now := time.Now()
	weekday := strings.ToLower(now.Weekday().String())
	hour := strconv.FormatInt(int64(now.Hour()), 10)

	// WHY DOES THIS WORK
	var price float64
	if result := Db.Raw("select ((pricing->?)::JSON->>CAST(? AS INT))::DECIMAL as price from spots where id=?", weekday, hour, spot.ID).Scan(&price); result.Error != nil {
		return 0
	}
	return price
}
