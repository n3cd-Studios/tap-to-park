package database

import (
	"time"
)

// Organization has many Spots, OrganizationID is the foreign key
// Organization has many Invites, OrganizationID is the foreign key
type Organization struct {
	ID      uint     `gorm:"primarykey" json:"-"`
	Name    string   `gorm:"not null;unique;" json:"name"`
	Spots   []Spot   `json:"spots"`
	Invites []Invite `json:"invites"`
}

type Invite struct {
	Code           string    `gorm:"primarykey;unique;default:upper(substr(md5(random()::text), 1, 10))" json:"code"`
	Expiration     time.Time `gorm:"not null;" json:"expiration"`
	OrganizationID uint      `gorm:"not null;" json:"organization"`
	CreatedByID    uint      `gorm:"not null;" json:"createdBy"`
	UsedByID       uint      `gorm:"" json:"usedBy"`
}

type UserRole = uint

const (
	USER UserRole = iota
	ADMIN
)

// User has many Sessions, UserID is the foreign key
type User struct {
	ID             uint      `gorm:"primarykey" json:"-"`
	Guid           string    `gorm:"not null;type:uuid;unique;default:gen_random_uuid()" json:"guid"`
	Email          string    `gorm:"not null;unique;" json:"email"`
	Role           UserRole  `gorm:"not null;default:0" json:"role"`
	PasswordHash   string    `gorm:"not null;default:''" json:"-"`
	Type           string    `gorm:"not null;default:'local'" json:"-"`
	ExternalID     string    `gorm:"not null;default:''" json:"-"`
	Sessions       []Session `json:"sessions"`
	OrganizationID uint      `json:"-"`
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
	UserID              uint      `gorm:"" json:"-"`
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
