package dto

type AvailableSeatsRequest struct {
	ScreeningID uint `json:"screening_id" binding:"required"`
}

type BookSeatsRequest struct {
	SeatIDs     []uint `json:"seat_ids" binding:"required"`
	ScreeningID uint   `json:"screening_id" binding:"required"`
	UserID      uint   `json:"user_id" binding:"required"`
}

type ReleaseSeatsRequest struct {
	BookingID uint `json:"booking_id" binding:"required"`
	UserID    uint `json:"user_id" binding:"required"`
}
