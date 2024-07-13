package handler

import (
	"ecoswap/api/token"
	pb "ecoswap/genproto/users"
	"ecoswap/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
// @Router /register [post]
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
// @Router /login [post]
func(h *Handler) Login(c *gin.Context){
	req := pb.UserLogin{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, err)
		h.Logger.Error(fmt.Sprintf("Login uchun ma'lumotlarni o'qishda xato: %v", err))
		return
	}

	user, err := h.UserRepo.GetUserByEmail(req.Email)
	if err != nil{
		c.JSON(http.StatusBadRequest, err)
		h.Logger.Error(fmt.Sprintf("Bazadan ma'lumotlarni o'qishda xato: %v", err))
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(user.Password)); err != nil{
		c.JSON(http.StatusBadRequest, err)
		h.Logger.Error(fmt.Sprintf("Parol xato: %v", err))
		return
	}

	token := token.GenerateJWT(&pb.GenerateUserJWT{
		Id: user.Id,
		Email: req.Email,
		Username: user.Username,
		FullName: user.FullName,
	})

	err = h.UserRepo.StoreRefreshToken(&model.RefreshToken{
		UserId: user.Id,
		Token: token.RefreshToken,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil{
		c.JSON(http.StatusInternalServerError, err)
		h.Logger.Error(fmt.Sprintf("Refresh token yaratilmadi: %v", err))
		return 
	}
	c.JSON(http.StatusOK, token)
}

func(h *Handler) LogOut(c *gin.Context){

}
