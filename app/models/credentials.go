package models

type Credentials struct {
	ID       uint
	Password string `gorm:"uniqueIndex;size:20"`
}
