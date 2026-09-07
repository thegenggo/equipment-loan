package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thegenggo/equipment-loan/api/internal/httperr"
	"github.com/thegenggo/equipment-loan/api/internal/middleware"
	"github.com/thegenggo/equipment-loan/api/internal/model"
	"github.com/thegenggo/equipment-loan/api/internal/service"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Name     string `json:"name" binding:"required,max=100"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.Write(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	user, err := h.auth.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			httperr.Write(c, http.StatusConflict, "email_taken", "this email is already registered")
			return
		}
		log.Printf("register: %v", err)
		httperr.Write(c, http.StatusInternalServerError, "internal_error", "could not create the account")
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.Write(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	signed, user, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			httperr.Write(c, http.StatusUnauthorized, "invalid_credentials", "email or password is incorrect")
			return
		}
		log.Printf("login: %v", err)
		httperr.Write(c, http.StatusInternalServerError, "internal_error", "could not sign in")
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: signed, User: toUserResponse(user)})
}

func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":   c.GetInt64(middleware.ContextUserID),
		"role": c.GetString(middleware.ContextRole),
	})
}

func toUserResponse(user *model.User) userResponse {
	return userResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}
}
