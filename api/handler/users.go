package handler

import (
	"ecoswap/api/token"
	"ecoswap/config"
	pb "ecoswap/genproto/users"
	"ecoswap/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body users.UserRegister true "User Register"
// @Success 200 {string} string "Muvaffaqiyatli ro'yxatdan o'tdingiz"
// @Failure 400 {object} model.Error "Xato"
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	req := pb.UserRegister{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		h.Logger.Error(fmt.Sprintf("Register uchun ma'lumotlarni o'qishda xato: %v", err))
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		h.Logger.Error(fmt.Sprintf("Passwordni hashlashda xato: %v", err))
		return
	}
	req.Password = string(hashPassword)

	err = h.UserRepo.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		h.Logger.Error(fmt.Sprintf("Bazaga yozishda xato: %v", err))
		return
	}
	c.JSON(http.StatusOK, "Muvaffaqiyatli ro'yxatdan o'tdingiz")
}

// @Summary Login user
// @Description user uchun login
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body users.UserLogin true "User Register"
// @Success 200 {string} users.Token
// @Failure 400 {object} model.Error "Xato"
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	req := pb.UserLogin{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		h.Logger.Error(fmt.Sprintf("Login uchun ma'lumotlarni o'qishda xato: %v", err))
		return
	}

	user, err := h.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		h.Logger.Error(fmt.Sprintf("Bazadan ma'lumotlarni o'qishda xato: %v", err))
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		h.Logger.Error(fmt.Sprintf("Parol xato: %v", err))
		return
	}

	token := token.GenerateJWT(&pb.GenerateUserJWT{
		Id:       user.Id,
		Email:    req.Email,
		Username: user.Username,
		FullName: user.FullName,
	})

	err = h.UserRepo.StoreRefreshToken(&model.RefreshToken{
		UserId:    user.Id,
		Token:     token.RefreshToken,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		h.Logger.Error(fmt.Sprintf("Refresh token yaratilmadi: %v", err))
		return
	}
	c.JSON(http.StatusOK, token)
}

// @Summary Reset password
// @Description Reset a user's password with the provided email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param email body users.Email true "User Email"
// @Success 200 {string} string "Password reset successful"
// @Failure 400 {object} string "Invalid request"
// @Failure 500 {object} string "Internal server error"
// @Router /auth/resetPass [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	req := pb.Email{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Ma'lumotlarni o'qishda xato: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := h.UserRepo.ResetPassword(&req)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Bazada ma'lumotlar topilmadi: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}
	c.JSON(http.StatusOK, resp)
}


// @Summary      Foydalanuvchi parolini yangilash
// @Description  Ushbu endpoint foydalanuvchi parolini yangilash uchun ishlatiladi.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        reset_password  body      model.ResetPassword  true  "Reset Password Payload"
// @Success      200             {object}  users.Status            "Successful Response"
// @Failure      400             {object}  users.Status            "Bad Request"
// @Router       /auth/updatePass [post]
func (h *Handler) UpdatePassword(c *gin.Context) {
	req := model.ResetPassword{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Ma'lumotlarni bodydan o'qishda xatolik: %v", err))
		c.JSON(http.StatusBadRequest, pb.Status{Status: false, Message: "Invalid request body"})
		return
	}

	resp, err := h.UserRepo.UpdatePassword(&req)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("Bazadan ma'lumotlarni olishda xato: %v", err))
		c.JSON(http.StatusBadRequest, pb.Status{Status: false, Message: "Failed to update password"})
		return
	}
	c.JSON(http.StatusOK, resp)
}


// @Summary      Foydalanuvchini tizimdan chiqarish
// @Description  Ushbu endpoint foydalanuvchini tizimdan chiqarish uchun ishlatiladi.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer <access_token>"
// @Success      200  {object}  map[string]string  "Logout successful"
// @Failure      400  {object}  map[string]string  "Invalid token or missing Authorization header"
// @Failure      500  {object}  map[string]string  "Failed to blacklist access token or cancel token"
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.Logger.Error("Token olinmadi")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Load().SIGNING_KEY), nil
	})

	if err != nil {
		h.Logger.Error("Token parsing error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"].(string)

		err := h.Redis.AddToBlacklist(c, accessToken, time.Hour*24)
		if err != nil {
			h.Logger.Error("Failed to blacklist access token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to blacklist access token"})
			return
		}

		err = h.UserRepo.CancelToken(userID)
		if err != nil {
			h.Logger.Error("Failed to cancel token: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	} else {
		h.Logger.Error("Invalid token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}
}

// @Summary      Foydalanuvchi tokenini yangilash
// @Description  Ushbu endpoint foydalanuvchi tokenini yangilash uchun ishlatiladi.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer <access_token>"
// @Success      200  {object}  map[string]string  "access_token"
// @Failure      400  {object}  map[string]string  "Invalid token or missing Authorization header"
// @Failure      500  {object}  map[string]string  "Failed to update token"
// @Router       /auth/updateToken [put]
func (h *Handler) UpdateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.Logger.Error("Token olinmadi")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := token.ExtraClaims(accessToken)
	if err != nil {
		h.Logger.Error("Token parsing error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	userID := claims["id"].(string)
	userEmail := claims["email"].(string)
	userUsername := claims["username"].(string)
	userFullName := claims["full_name"].(string)

	token := token.GenerateJWT(&pb.GenerateUserJWT{
		Id:       userID,
		Email:    userEmail,
		Username: userUsername,
		FullName: userFullName,
	})

	err = h.UserRepo.UpdateToken(&model.RefreshToken{
		UserId:    userID,
		Token:     token.RefreshToken,
		ExpiresAt: int64(time.Hour * 24),
	})
	if err != nil {
		h.Logger.Error("Yangi tokenni xotiraga yozishda xato: %v", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}
