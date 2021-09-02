package wallet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler interface {
	Create() gin.HandlerFunc
	Get() gin.HandlerFunc
}

type handler struct {
	s   Service
	log *zap.SugaredLogger
}

func NewHandler(s Service, log *zap.SugaredLogger) Handler {
	return &handler{
		s:   s,
		log: log,
	}
}

func (h *handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		wallet, err := h.s.Create()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, wallet)
	}
}

func (h *handler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		seed := c.Query("seed")
		if seed == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "seed is required",
			})
			return
		}

		i, err := strconv.ParseInt(c.Query("index"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		wallet, err := h.s.Get(seed, int(i))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, wallet)

	}
}
