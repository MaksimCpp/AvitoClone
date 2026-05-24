package handler

import (
	"errors"
	"net/http"
	"net/mail"

	errorresponse "github.com/MaksimCpp/AvitoClone/internal/delivery/http/error"
	"github.com/MaksimCpp/AvitoClone/internal/domain/user"
	jwtservice "github.com/MaksimCpp/AvitoClone/internal/infrastructure/jwt"
	userusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	registerUseCase userusecase.RegisterUserUseCase
	loginUseCase userusecase.LoginUserUseCase
	getMeUseCase userusecase.GetMeUseCase
	jwtService *jwtservice.JWTService
}

func NewUserHandler(
	registerUseCase userusecase.RegisterUserUseCase,
	loginUseCase userusecase.LoginUserUseCase,
	getMeUseCase userusecase.GetMeUseCase,
) *UserHandler {
	return &UserHandler{
		registerUseCase: registerUseCase,
		loginUseCase: loginUseCase,
		getMeUseCase: getMeUseCase,
	}
}

type RegisterAndLoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// Register godoc
//
// @Summary Register user
// @Description Register new user
// @Tags auth
// @Accept json
// @Produce json
//
// @Param request body RegisterAndLoginRequest true "Register request"
//
// @Success 201 {object} UserResponse
// @Failure 400 {object} errorresponse.ErrorResponse
//
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterAndLoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Password length < 6."})
		return
	}

	entity := user.User{
		Email: req.Email,
		Password: req.Password,
	}
	result, err := h.registerUseCase.Execute(c.Request.Context(), &entity)
	
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, errorresponse.ErrorResponse{Detail: err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	response := UserResponse{
		ID: result.ID,
		Email: result.Email,
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
//
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
//
// @Param request body RegisterAndLoginRequest true "Login request"
//
// @Success 200 {object} AccessToken
// @Failure 401 {object} errorresponse.ErrorResponse
//
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req RegisterAndLoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	entity := user.User{
		Email: req.Email,
		Password: req.Password,
	}
	token, err := h.loginUseCase.Execute(c.Request.Context(), &entity)

	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: err.Error()})
		} else if errors.Is(err, user.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, errorresponse.ErrorResponse{Detail: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, AccessTokenResponse{AccessToken: token})
}

// GetMe godoc
//
// @Summary Get current user
// @Description Returns current authorized user
// @Tags users
// @Accept json
// @Produce json
//
// @Security BearerAuth
//
// @Success 200 {object} UserResponse
// @Failure 401 {object} errorresponse.ErrorResponse
//
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, errorresponse.ErrorResponse{Detail: "Unauthorized."})
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, errorresponse.ErrorResponse{Detail: "Invalid user id."})
		return
	}
	result, err := h.getMeUseCase.Execute(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorresponse.ErrorResponse{Detail: err.Error()})
		return
	}

	response := UserResponse{
		ID: result.ID,
		Email: result.Email,
	}
	c.JSON(http.StatusOK, response)
}
