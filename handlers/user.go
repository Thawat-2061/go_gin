package handlers

import (
	"net/http"

	"go-gin/dto"
	"go-gin/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) Register(c *gin.Context) {
	// 1. กำหนดโครงสร้าง request
	var request struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=8"`
		Email    string `json:"email" binding:"required,email"`
	}

	// 2. ตรวจสอบและ bind ข้อมูล
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"code":   "INVALID_INPUT",
			"error":  err.Error(),
		})
		return
	}

	// 3. ตรวจสอบ username และ email ที่มีอยู่แล้ว
	var existingUser models.User
	if err := h.db.Where("(username = ? OR email = ?) AND deleted_at IS NULL", request.Username, request.Email).First(&existingUser).Error; err == nil {
		status := http.StatusConflict
		errorCode := "USER_EXISTS"
		message := "User already exists"

		if existingUser.Email == request.Email {
			errorCode = "EMAIL_EXISTS"
			message = "Email already registered"
		} else if existingUser.Username == request.Username {
			errorCode = "USERNAME_EXISTS"
			message = "Username already taken"
		}

		c.JSON(status, gin.H{
			"status":  "error",
			"code":    errorCode,
			"message": message,
		})
		return
	}

	// 4. เข้ารหัส password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "PASSWORD_HASH_FAILED",
			"message": "Could not hash password",
		})
		return
	}

	// 5. สร้างผู้ใช้ใหม่
	newUser := models.User{
		Username: request.Username,
		Password: string(hashedPassword),
		Email:    request.Email,
	}

	if err := h.db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"code":    "USER_CREATION_FAILED",
			"message": "Could not create user",
			"error":   err.Error(),
		})
		return
	}

	// 6. ตอบกลับผลสำเร็จ
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"data": gin.H{
			"user": gin.H{
				"id":         newUser.ID,
				"username":   newUser.Username,
				"email":      newUser.Email,
				"created_at": newUser.CreatedAt,
			},
		},
	})
}
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง"})
		return
	}

	// ในที่นี้เราคืนค่า user_id ง่ายๆ
	// ในทางปฏิบัติควรสร้างและคืนค่า JWT token
	c.JSON(http.StatusOK, gin.H{
		"message": "ล็อกอินสำเร็จ",
		"user_id": user.ID,
	})
}
