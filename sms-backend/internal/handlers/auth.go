package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sms-backend/internal/helpers"
	"sms-backend/internal/services"
)

type AuthHandler struct{ service *services.AuthService }

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}
func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.service.RevokeUserTokens(c.GetUint("userID")); err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired session")
		return
	}
	helpers.SuccessResponse(c, http.StatusOK, "logged out")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}
	if c.ShouldBindJSON(&req) != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "university ID and password are required")
		return
	}
	user, token, err := h.service.Login(req.Identifier, req.Password)
	if err != nil {
		helpers.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	helpers.SuccessResponse(c, http.StatusOK, gin.H{"token": token, "user": user})
}
func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.service.GetUser(c.GetUint("userID"))
	if err != nil {
		helpers.ErrorResponse(c, 404, "user not found")
		return
	}
	helpers.SuccessResponse(c, 200, user)
}
func (h *AuthHandler) UpdateMe(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email" binding:"omitempty,email"`
	}
	if c.ShouldBindJSON(&req) != nil {
		helpers.ErrorResponse(c, 400, "invalid account details")
		return
	}
	user, err := h.service.UpdateUser(c.GetUint("userID"), req.Name, req.Email)
	if err != nil {
		helpers.ErrorResponse(c, 400, err.Error())
		return
	}
	helpers.SuccessResponse(c, 200, user)
}
func (h *AuthHandler) Password(c *gin.Context) {
	var req struct {
		Current string `json:"currentPassword" binding:"required"`
		Next    string `json:"newPassword" binding:"required"`
	}
	if c.ShouldBindJSON(&req) != nil {
		helpers.ErrorResponse(c, 400, "both passwords are required")
		return
	}
	if err := h.service.ChangePassword(c.GetUint("userID"), req.Current, req.Next); err != nil {
		helpers.ErrorResponse(c, 400, err.Error())
		return
	}
	helpers.SuccessResponse(c, 200, "password updated")
}
func (h *AuthHandler) DeleteMe(c *gin.Context) {
	if c.GetString("role") != "admin" {
		helpers.ErrorResponse(c, http.StatusForbidden, "student and professor accounts can only be deleted by an administrator")
		return
	}
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if c.ShouldBindJSON(&req) != nil {
		helpers.ErrorResponse(c, 400, "password is required")
		return
	}
	if err := h.service.DeleteUser(c.GetUint("userID"), req.Password); err != nil {
		helpers.ErrorResponse(c, 400, err.Error())
		return
	}
	helpers.SuccessResponse(c, 200, "account deleted")
}
