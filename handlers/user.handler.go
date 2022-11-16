package handler

import (
	"cyberzell.com/seguros/models"
	"cyberzell.com/seguros/models/apperrors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	
)

func (h *Handler) GetUserById(c *gin.Context) {
	routeId := c.Param("id")

	userId, _ := strconv.Atoi(routeId)

	user, err := h.userService.GetUserById(userId)

	if err != nil {
		log.Printf("Unable to find user for id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", routeId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

type loginPayload struct {
	Email		string	`json:"email"`
	Password	string	`json:"password"`
}

func (r loginPayload) validate() error {
	return validation.ValidateStruct(&r,
			validation.Field(&r.Email, validation.Required, is.EmailFormat),
			validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
		)
}

func (r *loginPayload) sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Login(c *gin.Context) {
	var request loginPayload

	if ok := bindData(c, &request); !ok {
		return
	}

	request.sanitize()
	user, err := h.userService.Login(request.Email, request.Password)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error":	err,
			"status":	http.StatusBadRequest,
		})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:		strconv.Itoa(int(user.Id)),
		ExpiresAt:	time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON((http.StatusInternalServerError), gin.H{
			"status_code":		http.StatusInternalServerError,
			"message":			"Incorrect password",
		})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:		"jwt",
		Value:		token,
		MaxAge:		0,
		Path:		"/",
	})

	c.JSON(http.StatusOK, user)
}

type registerPayload struct {
	Email		string	`json:"email"`
	Username		string	`json:"username"`
	Password		string	`json:"password"`
}

func (r registerPayload) validate() error {
	return validation.ValidateStruct(&r,
			validation.Field(&r.Email, validation.Required, is.EmailFormat),
			validation.Field(&r.Username, validation.Required, validation.Length(3, 30)),
			validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *registerPayload) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Register(c *gin.Context) {
	var request registerPayload

	if ok := bindData(c, &request); !ok {
		return
	}

	request.sanitize()
	registerUserPayload := &models.User{
		Username:	request.Username,
		Email:		request.Email,
		Password:	request.Password,
	}
	user, err := h.userService.Register(registerUserPayload)

	if err != nil {
		if err.Error() == apperrors.NewBadRequest(apperrors.DuplicateEmail).Error() {
			toFieldErrorResponse(c, "Email", apperrors.DuplicateEmail)
			return
		}
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:		"jwt",
		Value:		"",
		MaxAge:		0,
		Path:		"/",
		Expires:	time.Now().Add(-time.Hour),
	})

	c.JSON((http.StatusOK), gin.H{
		"status_code":	http.StatusOK,
		"message":		"Logged out",
	})
}