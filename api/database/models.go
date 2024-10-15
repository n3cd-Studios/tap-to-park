package database

import (
	"time"

	"github.com/jackc/pgx/pgtype"
)

// Organization has many Spots, OrganizationID is the foreign key
// Organization has many Admins, OrganizationID is the foreign key
// Organization has many Invites, OrganizationID is the foreign key
type Organization struct {
	ID      uint     `gorm:"primarykey" json:"-"`
	Name    string   `gorm:"not null;unique;" json:"name"`
	Spots   []Spot   `json:"spots"`
	Invites []Invite `json:"invites"`
	Admins  []User   `json:"-"`
}

type Invite struct {
	Code           string    `gorm:"primarykey;unique;default:upper(substr(md5(random()::text), 1, 10))" json:"code"`
	Expiration     time.Time `gorm:"not null;" json:"expiration"`
	OrganizationID uint      `gorm:"not null;" json:"organization"`
	CreatedByID    uint      `gorm:"not null;" json:"createdBy"`
	UsedByID       uint      `gorm:"" json:"usedBy"`
}

type User struct {
	ID             uint   `gorm:"primarykey" json:"-"`
	Guid           string `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Email          string `gorm:"not null;unique;" json:"email"`
	PasswordHash   string `gorm:"not null;" json:"-"`
	OrganizationID uint   `gorm:"not null;" json:"-"`
}

// Spot has many Reservations, SpotID is the foreign key
type Spot struct {
	ID             uint          `gorm:"primarykey" json:"-"`
	Guid           string        `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Name           string        `gorm:"not null;" json:"name"`
	Coords         Coordinates   `gorm:"type:Point;index:coords_gist_idx,type:gist" json:"coords"`
	Handicap       bool          `gorm:"not null;" json:"handicap"`
	OrganizationID uint          `gorm:"not null;" json:"organization"`
	Reservations   []Reservation `json:"reservations"`
	Table          pgtype.JSONB  `gorm:"type:jsonb;not null;serializer:json" json:"table"`
}

type Reservation struct {
	ID            uint      `gorm:"primarykey" json:"-"`
	Guid          string    `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Start         time.Time `gorm:"not null;" json:"start"`
	End           time.Time `gorm:"not null;" json:"end"`
	Cost          float64   `gorm:"not null;" json:"cost"`
	TransactionID string    `gorm:"not null;" json:"-"`
	SpotID        uint      `gorm:"not null;" json:"-"`
}
