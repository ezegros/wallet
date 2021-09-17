package wallet

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ezegrosfeld/wallet/users/pkg/security"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Create() gin.HandlerFunc
	GetAddress() gin.HandlerFunc
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

// Create returns a gin handler function which creates a new wallet for a user
func (h *handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := c.MustGet("user").(*security.JWTData)

		wallet, err := h.s.Create(c, data.UserID)
		if err != nil {
			if errors.Is(err, ErrConflict) {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, wallet)
	}
}

// GetAddress returns a gin handler function which gets an address for a user
func (h *handler) GetAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := c.MustGet("user").(*security.JWTData)

		idx := c.Query("index")

		/* if idx == "" {
			addr, err := h.s.Get(c, data.UserID)
			if err != nil {
				if errors.Is(err, ErrNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, addr)
			return
		} */

		i, err := strconv.ParseInt(idx, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		address, err := h.s.GetAddress(c, data.UserID, int(i))
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, address)
	}
}
