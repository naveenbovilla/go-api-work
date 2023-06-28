package model

import (
	"time"

	"gorm.io/gorm"
)

type ExtraHours struct {
	gorm.Model
	ServiceID      uint      `json:"service_id"`
	AcceptDate     time.Time `json:"accept_date"`
	RefuseDate     time.Time `json:"refuse_date"`
	WaitingDate    time.Time `json:"waiting_date"`
	HoursRequested uint      `json:"hours_requested"`
	TypeOfWork     string    `json:"type_of_work"`
	Notes          string    `json:"notes"`
	Status         string    `json:"status"`
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&ExtraHours{})
	db.Exec("ALTER TABLE extra_hours ADD PRIMARY KEY (service_id);")
	if err != nil {
		// Handle the migration error appropriately
		panic(err)
	}
}

func GetNextServiceID(db *gorm.DB) (uint, error) {
	var maxID uint
	result := db.Table("extra_hours").Select("COALESCE(MAX(service_id), 0)").Scan(&maxID)
	if result.Error != nil {
		return 0, result.Error
	}
	return maxID + 1, nil
}
