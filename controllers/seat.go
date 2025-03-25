package controllers

import (
	"go-gin/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SeatController(router *gin.Engine, db *gorm.DB) {
	handler := handlers.NewSeatHandler(db)

	seatGroup := router.Group("/seats")
	{
		// แก้ไขบรรทัดนี้ - เปลี่ยนจาก POST เป็น GET หากต้องการใช้ GET method
		seatGroup.POST("/available", handler.GetAvailableSeats) // หรือใช้ GET ตามที่ต้องการ
		seatGroup.POST("/book", handler.BookSeats)
		seatGroup.POST("/release", handler.ReleaseSeats)
	}

	bookingGroup := router.Group("/bookings")
	{
		bookingGroup.GET("/:user_id", handler.GetUserBookings)
	}
}
