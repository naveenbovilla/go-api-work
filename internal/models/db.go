package model

import (
	"time"

	"gorm.io/gorm"
)

type ExtraHours struct {
	gorm.Model
	// TODO adds quantity extra hours and other fields that you think are needed.
	ServiceID   uint
	AcceptDate  time.Time
	RefuseDate  time.Time
	WaitingDate time.Time
	Status      string
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&ExtraHours{})
}
