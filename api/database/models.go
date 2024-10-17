package database

import (
	"time"
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

// User has many Sessions, UserID is the foreign key
type User struct {
	ID             uint      `gorm:"primarykey" json:"-"`
	Guid           string    `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Email          string    `gorm:"not null;unique;" json:"email"`
	PasswordHash   string    `gorm:"not null;" json:"-"`
	Sessions       []Session `json:"sessions"`
	OrganizationID uint      `gorm:"not null;" json:"-"`
}

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

type Reservation struct {
	ID                  uint      `gorm:"primarykey" json:"-"`
	Guid                string    `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Start               time.Time `gorm:"not null;" json:"start"`
	End                 time.Time `gorm:"not null;" json:"end"`
	Price               float64   `gorm:"not null;" json:"price"`
	Email               string    `gorm:"not null;" json:"email"`
	StripeTransactionID string    `gorm:"not null;unique;" json:"-"`
	SpotID              uint      `gorm:"not null;" json:"-"`
	// Transactions []Transaction `gorm:"not null;" json:"transactions"`
}

type Session struct {
	ID       uint      `gorm:"primarykey" json:"-"`
	Guid     string    `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	IP       string    `gorm:"not null;" json:"ip"`
	Device   string    `gorm:"not null;" json:"device"`
	Expires  time.Time `gorm:"not null;" json:"expires"`
	LastUsed time.Time `gorm:"not null;" json:"lastUsed"`
	UserID   uint      `gorm:"not null;" json:"-"`
}
