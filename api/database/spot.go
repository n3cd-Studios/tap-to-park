package database

import (
	"strconv"
	"strings"
	"time"
)

// Spot has many Reservations, SpotID is the foreign key
type Spot struct {
	ID             uint          `gorm:"primarykey" json:"-"`
	Guid           string        `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Name           string        `gorm:"not null;" json:"name"`
	Coords         Coordinates   `gorm:"type:Point;index:coords_gist_idx,type:gist" json:"coords"`
	Handicap       bool          `gorm:"not null;" json:"handicap"`
	OrganizationID uint          `gorm:"not null;" json:"organization"`
	Pricing        Pricing       `gorm:"type:json;not null;default:'{\"sunday\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"monday\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"tuesday\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"wednesday\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"thursday\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"friday\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"saturday\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]}';" json:"table"`
	Reservations   []Reservation `json:"reservations"`
}

func (spot *Spot) GetReservation() *Reservation {

	reservation := &Reservation{}
	if result := Db.Where("spot_id = ?", spot.ID).Where("\"end\" > ?", time.Now()).First(reservation); result.Error != nil {
		return nil
	}

	return reservation
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
