package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/ezegrosfeld/wallet/users/pkg/security"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Create() gin.HandlerFunc
}

type handler struct {
	s Service
}

// NewHandler creates a new Handler
func NewHandler(s Service) Handler {
	return &handler{
		s: s,
	}
}

// Create creates returns handlerFunc which creates a new user
func (h *handler) Create() gin.HandlerFunc {
	// Expected request
	type request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		// Bind the request
		r := new(request)
		err := c.ShouldBindJSON(&r)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create the user using the service
		user, err := h.s.Create(c, r.Username, r.Password)
		if err != nil {
			if errors.Is(err, ErrConflict) {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token, _ := security.CreateJWT(user.ID)

		c.SetCookie("token", token, int(time.Hour)*24*7, "/", "", false, true)

		c.JSON(http.StatusCreated, user)
	}
}
