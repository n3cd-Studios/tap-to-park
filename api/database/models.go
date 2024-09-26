package database

import (
	"time"
)

// Organization has many Spots, OrganizationID is the foreign key
// Organization has many Admins, OrganizationID is the foreign key
type Organization struct {
	ID     uint   `gorm:"primarykey" json:"-"`
	Name   string `gorm:"not null;unique;" json:"name"`
	Spots  []Spot `json:"spots"`
	Admins []User `json:"-"`
}

type User struct {
	ID             uint   `gorm:"primarykey" json:"-"`
	Guid           string `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Email          string `gorm:"not null;unique;" json:"email"`
	PasswordHash   string `gorm:"not null;" json:"-"`
	OrganizationID uint   `gorm:"not null;" json:"-"`
}

type Reservation struct {
	ID            uint      `gorm:"primarykey" json:"-"`
	Guid          string    `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Start         time.Time `gorm:"not null;" json:"start"`
	End           time.Time `gorm:"not null;" json:"end"`
	SpotID        uint      `gorm:"not null;" json:"id"`
	TransactionID string    `gorm:"not null;" json:"-"`
	CostPerHour   uint      `gorm:"not null;" json:"costPerHour"`
}

type Spot struct {
	ID             uint        `gorm:"primarykey" json:"-"`
	Guid           string      `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Name           string      `gorm:"not null;" json:"name"`
	Coords         Coordinates `gorm:"type:Point;index:coords_gist_idx,type:gist" json:"coords"`
	Handicap       bool        `gorm:"not null;" json:"handicap"`
	OrganizationID uint        `gorm:"not null;" json:"organization"`
	Reservations   []Reservation
}

type Error struct {
	Message string `json:"message"`
}
