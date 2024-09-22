package database

import (
	"time"

	"gorm.io/gorm"
)

// Organization has many Spots, OrganizationID is the foreign key
// Organization has many Admins, OrganizationID is the foreign key
type Organization struct {
	gorm.Model
	Name   string
	Spots  []Spot
	Admins []User
}

type User struct {
	gorm.Model
	Email          string
	PasswordHash   string
	OrganizationID uint
}

type Reservation struct {
	gorm.Model
	Start         time.Time
	End           time.Time
	SpotID        uint
	TransactionID string
}

type Spot struct {
	gorm.Model
	Coords         Point `gorm:"type:Point; index:coords_gist_idx, type:gist"`
	Handicap       bool
	OrganizationID uint
}
