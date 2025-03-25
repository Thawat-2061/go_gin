package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"go-gin/dto"
	"go-gin/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SeatHandler struct {
	db *gorm.DB
}

func NewSeatHandler(db *gorm.DB) *SeatHandler {
	return &SeatHandler{db: db}
}

// GetAvailableSeats ดึงที่นั่งว่าง
func (h *SeatHandler) GetAvailableSeats(c *gin.Context) {
	var req dto.AvailableSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var seats []models.Seat
	if err := h.db.Where("screening_id = ? AND status = ?", req.ScreeningID, "available").Find(&seats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลที่นั่งได้"})
		return
	}

	c.JSON(http.StatusOK, seats)
}

// BookSeats จองที่นั่ง
func (h *SeatHandler) BookSeats(c *gin.Context) {
	var req dto.BookSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบว่ามีผู้ใช้นี้จริงหรือไม่
	var user models.User
	if err := h.db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่พบผู้ใช้"})
		return
	}

	// ตรวจสอบที่นั่งว่างทั้งหมดก่อน
	var availableSeats int64
	h.db.Model(&models.Seat{}).
		Where("id IN ? AND status = ?", req.SeatIDs, "available").
		Count(&availableSeats)

	if availableSeats != int64(len(req.SeatIDs)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ที่นั่งบางส่วนถูกจองแล้ว"})
		return
	}

	// สร้างการจอง
	booking := models.Booking{
		UserID:      req.UserID,
		SeatIDs:     strings.Trim(strings.Join(strings.Fields(fmt.Sprint(req.SeatIDs)), ","), "[]"),
		ScreeningID: req.ScreeningID,
	}

	// บันทึกการจองและอัปเดตที่นั่ง
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Seat{}).
			Where("id IN ?", req.SeatIDs).
			Updates(map[string]interface{}{
				"status":     "booked",
				"booking_id": booking.ID,
			}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "การจองไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "จองสำเร็จ",
		"booking_id": booking.ID,
		"user_id":    user.ID,
		"username":   user.Username,
	})
}

// ReleaseSeats ปล่อยที่นั่ง
func (h *SeatHandler) ReleaseSeats(c *gin.Context) {
	var req dto.ReleaseSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบสิทธิ์ผู้ใช้ (ควรใช้ middleware ในทางปฏิบัติ)
	var booking models.Booking
	if err := h.db.First(&booking, req.BookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบการจองนี้"})
		return
	}

	if booking.UserID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "ไม่มีสิทธิ์ปล่อยที่นั่งนี้"})
		return
	}

	// แปลง seat_ids จาก string เป็น slice
	seatIDs := strings.Split(booking.SeatIDs, ",")

	// ปล่อยที่นั่ง
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// อัปเดตสถานะที่นั่ง
		if err := tx.Model(&models.Seat{}).
			Where("id IN ?", seatIDs).
			Updates(map[string]interface{}{
				"status":     "available",
				"booking_id": nil,
			}).Error; err != nil {
			return err
		}

		// ลบการจอง
		if err := tx.Delete(&booking).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ปล่อยที่นั่งไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ปล่อยที่นั่งสำเร็จ"})
}

// GetUserBookings ดึงประวัติการจองของผู้ใช้
func (h *SeatHandler) GetUserBookings(c *gin.Context) {
	userID := c.Param("user_id")

	var bookings []models.Booking
	if err := h.db.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงประวัติการจองได้"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
