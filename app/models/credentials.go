package models

// Credentials represents user credentials.
type Credentials struct {
	ID       uint
	Password string `gorm:"uniqueIndex;size:256"`
}
