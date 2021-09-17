package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ezegrosfeld/wallet/users/internal/domain"
	"github.com/ezegrosfeld/wallet/users/pkg/security"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedService struct {
	mock.Mock
}

func (m *mockedService) Create(ctx context.Context, userID string) (*domain.Wallet, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}
func (m *mockedService) Get(ctx context.Context, userID string) (*domain.Wallet, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}
func (m *mockedService) GetAddress(ctx context.Context, userID string, index int) (string, error) {
	args := m.Called(ctx, userID, index)
	return args.String(0), args.Error(1)
}

func createServer(s *mockedService) *gin.Engine {
	h := NewHandler(s)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("user", &security.JWTData{
			UserID: "1",
		})
	})

	router.POST("/users/wallet", h.Create())
	router.GET("/users/wallet", h.GetAddress())

	return router
}

func createRequest(method, endpoint, body string) (*http.Request, *httptest.ResponseRecorder) {
	r := httptest.NewRequest(method, endpoint, strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	rw := httptest.NewRecorder()

	return r, rw
}

func TestHandlerCreate(t *testing.T) {
	s := &mockedService{}
	s.On("Create", mock.Anything, "1").Return(&domain.Wallet{ID: "123", UserID: "1", Seed: "seed"}, nil)

	router := createServer(s)

	r, rw := createRequest("POST", "/users/wallet", "")
	router.ServeHTTP(rw, r)

	type response struct {
		domain.Wallet
	}

	var res response
	err := json.Unmarshal(rw.Body.Bytes(), &res)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rw.Code)
	assert.Equal(t, "123", res.ID)
	assert.Equal(t, "1", res.UserID)
	assert.Equal(t, "seed", res.Seed)
}

func TestHandlerCreateWithConflict(t *testing.T) {
	s := &mockedService{}
	s.On("Create", mock.Anything, "1").Return(&domain.Wallet{}, ErrConflict)

	router := createServer(s)

	r, rw := createRequest("POST", "/users/wallet", "")
	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusConflict, rw.Code)
}

func TestHandlerCreateWithInternalServerError(t *testing.T) {
	s := &mockedService{}
	s.On("Create", mock.Anything, "1").Return(nil, fmt.Errorf("error"))

	router := createServer(s)

	r, rw := createRequest("POST", "/users/wallet", "")
	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
}

/* func TestHandlerGet(t *testing.T) {
	s := &mockedService{}
	s.On("Get", mock.Anything, "1").Return(&domain.Wallet{ID: "123", UserID: "1", Seed: "seed"}, nil)

	router := createServer(s)

	r, rw := createRequest("GET", "/users/wallet", "")
	router.ServeHTTP(rw, r)

	type response struct {
		domain.Wallet
	}

	var res response
	err := json.Unmarshal(rw.Body.Bytes(), &res)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "123", res.ID)
	assert.Equal(t, "1", res.UserID)
	assert.Equal(t, "seed", res.Seed)
}

func TestHandlerGetWithNotFound(t *testing.T) {
	s := &mockedService{}
	s.On("Get", mock.Anything, "1").Return(&domain.Wallet{}, ErrNotFound)

	router := createServer(s)

	r, rw := createRequest("GET", "/users/wallet", "")
	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusNotFound, rw.Code)
}

func TestHandlerGetWithInternalServerError(t *testing.T) {
	s := &mockedService{}
	s.On("Get", mock.Anything, "1").Return(nil, fmt.Errorf("error"))

	router := createServer(s)

	r, rw := createRequest("GET", "/users/wallet", "")
	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
} */

func TestHandlerGet(t *testing.T) {
	s := &mockedService{}

	s.On("GetAddress", mock.Anything, "1", 0).Return("supercomplexaddress", nil)

	router := createServer(s)

	r, rw := createRequest("GET", "/users/wallet?index=0", "")
	router.ServeHTTP(rw, r)

	type response string

	var res response
	err := json.Unmarshal(rw.Body.Bytes(), &res)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "supercomplexaddress", string(res))
}

func TestHandlerGetWithProvidedIndex1(t *testing.T) {
	s := &mockedService{}

	s.On("GetAddress", mock.Anything, "1", 1).Return("supercomplexaddress", nil)

	router := createServer(s)

	r, rw := createRequest("GET", "/users/wallet?index=1", "")
	router.ServeHTTP(rw, r)

	type response string

	var res response
	err := json.Unmarshal(rw.Body.Bytes(), &res)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Equal(t, "supercomplexaddress", string(res))
}

func TestHandlerGetWithNotFound(t *testing.T) {
	s := &mockedService{}
	s.On("GetAddress", mock.Anything, "1", 0).Return("", ErrNotFound)

	router := createServer(s)

	r, rw := createRequest("GET", "/users/wallet?index=0", "")
	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusNotFound, rw.Code)
}

func TestHandlerGetWithInternalServerError(t *testing.T) {
	s := &mockedService{}
	s.On("GetAddress", mock.Anything, "1", 0).Return("", fmt.Errorf("error"))

	router := createServer(s)

	r, rw := createRequest("GET", "/users/wallet?index=0", "")
	router.ServeHTTP(rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
}
