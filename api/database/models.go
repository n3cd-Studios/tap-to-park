package database

import (
	"time"

	"gorm.io/gorm"
)

// Organization has many Spots, OrganizationID is the foreign key
// Organization has many Admins, OrganizationID is the foreign key
type Organization struct {
	gorm.Model
	Name   string `gorm:"not null;unique;" json:"name"`
	Spots  []Spot
	Admins []User
}

type User struct {
	gorm.Model     `json:"-"`
	UniqueID       string `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"uuid"`
	Email          string `gorm:"not null;unique;" json:"email"`
	PasswordHash   string `gorm:"not null;" json:"-"`
	OrganizationID uint   `gorm:"not null;" json:"-"`
}

type Reservation struct {
	gorm.Model
	Start         time.Time `gorm:"not null;" json:"start"`
	End           time.Time `gorm:"not null;" json:"end"`
	SpotID        uint      `gorm:"not null;" json:"id"`
	TransactionID string    `gorm:"not null;" json:"-"`
}

type Spot struct {
	gorm.Model
	Name           string      `gorm:"not null;" json:"name"`
	Coords         Coordinates `gorm:"type:Point;index:coords_gist_idx,type:gist" json:"coords"`
	Handicap       bool        `gorm:"not null;" json:"handicap"`
	OrganizationID uint        `gorm:"not null;" json:"organization"`
}

type Error struct {
	Message string `json:"message"`
}
