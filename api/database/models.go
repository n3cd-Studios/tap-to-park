package database

import (
	"time"

	"gorm.io/gorm"
)

// Organization has many Spots, OrganizationID is the foreign key
// Organization has many Admins, OrganizationID is the foreign key
type Organization struct {
	gorm.Model
	Name   string `gorm:"not null;unique;"`
	Spots  []Spot
	Admins []User
}

type User struct {
	gorm.Model
	UniqueID       string `gorm:"not null;type:uuid;unique;default:gen_random_uuid()"`
	Email          string `gorm:"not null;unique;"`
	PasswordHash   string `gorm:"not null;"`
	OrganizationID uint   `gorm:"not null;"`
}

type Reservation struct {
	gorm.Model
	Start         time.Time `gorm:"not null;"`
	End           time.Time `gorm:"not null;"`
	SpotID        uint      `gorm:"not null;"`
	TransactionID string    `gorm:"not null;"`
}

type Spot struct {
	gorm.Model
	Coords         Point `gorm:"type:Point;index:coords_gist_idx,type:gist"`
	Handicap       bool  `gorm:"not null;"`
	OrganizationID uint  `gorm:"not null;"`
}

type Error struct {
	Message string `json:"message"`
}
