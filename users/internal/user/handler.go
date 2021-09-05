package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Create() gin.HandlerFunc
}

type handler struct {
	s Service
}

func NewHandler(s Service) Handler {
	return &handler{
		s: s,
	}
}

func (h *handler) Create() gin.HandlerFunc {
	type request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		r := new(request)
		err := c.ShouldBindJSON(&r)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := h.s.Create(c, r.Username, r.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}
