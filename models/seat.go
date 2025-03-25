package models

import "gorm.io/gorm"

type Seat struct {
	gorm.Model
	ScreenID    uint   `json:"screen_id"`
	Row         string `json:"row"`
	Number      int    `json:"number"`
	Status      string `json:"status" gorm:"default:available"` // available, booked
	ScreeningID uint   `json:"screening_id"`
}

type Booking struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	SeatIDs     string `json:"seat_ids"`
	ScreeningID uint   `json:"screening_id"`
	User        User   `gorm:"foreignKey:UserID"` // เพิ่มความสัมพันธ์
}
