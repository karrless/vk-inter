package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"
	"vk-inter/internal/transport/rest/interfaces"
	"vk-inter/pkg/errs"
	"vk-inter/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthController struct {
	ctx     *context.Context
	service interfaces.AuthService
}

func NewAuthController(ctx *context.Context, authService interfaces.AuthService) *AuthController {
	return &AuthController{
		ctx:     ctx,
		service: authService,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Login     string             `json:"login"`
	CreatedAt time.Time          `bson:"created_at"`
}

// @Summary	Sign up endpoint
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		request	body	SignUpRequest	true	"Login and password"
// @Success	201		{object}	SignUpResponse
// @Failure	400		{object}	ErrorResponse
// @Failure	409		{object}	ErrorResponse
// @Router		/auth/signup [post]
func (ac *AuthController) SignUp(c *gin.Context) {
	if isAuth, exists := c.Get("isAuthenticated"); exists && isAuth.(bool) {
		resp := ErrorResponse{
			Error:   errs.ErrUserAlreadyAuth.Error(),
			Message: "Please logout first to signup",
		}
		c.JSON(http.StatusConflict, resp)
		return
	}

	var req SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Login) < 3 || len(req.Login) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidLoginFormat})
		return
	}
	user, err := ac.service.SignUp(*ac.ctx, req.Login, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if _, ok := err.(utils.PasswordError); ok {
			status = http.StatusBadRequest
		}
		if errors.Is(err, errs.ErrInvalidLoginFormat) || errors.Is(err, errs.ErrInvalidPasswordLength) {
			status = http.StatusBadRequest
		}
		if errors.Is(err, errs.ErrUserAlreadyExsist) {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
	}
	resp := SignUpResponse{
		ID:        user.ID,
		Login:     user.Login,
		CreatedAt: user.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

type LogInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LogInResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expresIn"`
}

// @Summary	Sign up endpoint
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		request	body	SignUpRequest	true	"Login and password"
// @Success	201		{object}	SignUpResponse
// @Failure	400		{object}	ErrorResponse
// @Failure	403		{object}	ErrorResponse
// @Failure	409		{object}	ErrorResponse
// @Router		/auth/login [post]
func (ac *AuthController) LogIn(c *gin.Context) {
	if isAuth, exists := c.Get("isAuthenticated"); exists && isAuth.(bool) {
		tokenSubject, _ := c.Get("id")
		if tokenSubject != nil { // если нет id, то токен изначально плохой, потому пропускаем это
			user, err := ac.service.GetUserById(*ac.ctx, tokenSubject.(string))
			if err == nil {
				resp := ErrorResponse{Message: "Please logout first to login"}
				status := 0

				if tokenSubject.(string) != user.Login {
					resp.Error = errs.ErrAlreadyLoggedIn.Error()
					status = http.StatusForbidden
				} else {
					resp.Error = errs.ErrUserAlreadyAuth.Error()
					status = http.StatusConflict
				}
				c.JSON(status, resp)
				return
			}
		}
	}

	var req LogInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, expiresIn, err := ac.service.LogIn(*ac.ctx, req.Login, req.Password)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, errs.ErrWrongPasswordOrLogin) {
			status = http.StatusUnauthorized // Не 404, потому чтоб не раскрывать, по какой причине невозможно залогиниться
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	resp := LogInResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	}
	c.JSON(http.StatusOK, resp)
}
