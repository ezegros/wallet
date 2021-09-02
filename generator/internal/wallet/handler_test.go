package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ezegrosfeld/wallet/generator/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockedService struct {
	mock.Mock
}

func (m *mockedService) Create() (domain.Wallet, error) {
	args := m.Called()
	return args.Get(0).(domain.Wallet), args.Error(1)
}

func (m *mockedService) Get(seedString string, index int) (domain.Wallet, error) {
	args := m.Called(seedString, index)
	return args.Get(0).(domain.Wallet), args.Error(1)
}

// create a mocked service for testing
func createMockedService(s *mockedService) *gin.Engine {
	handler := NewHandler(s, &zap.SugaredLogger{})

	router := gin.Default()

	router.POST("/wallet", handler.Create())
	router.GET("/wallet", handler.Get())

	return router
}

func createRequest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

// Test Create
func TestCreate(t *testing.T) {
	type response struct {
		Seed    string `json:"seed"`
		Address string `json:"address"`
	}

	s := new(mockedService)
	s.On("Create").Return(domain.Wallet{
		Seed:    "seed",
		Address: "address",
	}, nil)

	router := createMockedService(s)
	req, rw := createRequest("POST", "/wallet", "")

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusCreated, rw.Code)

	res := response{}

	err := json.Unmarshal(rw.Body.Bytes(), &res)
	assert.Nil(t, err)

	assert.Equal(t, "seed", res.Seed)
	assert.Equal(t, "address", res.Address)
}

// Test create with error
func TestCreateError(t *testing.T) {
	s := new(mockedService)
	s.On("Create").Return(domain.Wallet{}, fmt.Errorf("An error ocurred while creating the wallet"))

	router := createMockedService(s)
	req, rw := createRequest("POST", "/wallet", "")

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
}

// Test Get
func TestGet(t *testing.T) {
	type response struct {
		Seed    string `json:"seed"`
		Address string `json:"address"`
	}

	s := new(mockedService)
	s.On("Get", "seed", 0).Return(domain.Wallet{
		Seed:    "seed",
		Address: "address",
	}, nil)

	router := createMockedService(s)
	req, rw := createRequest("GET", "/wallet?seed=seed&index=0", "")

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusOK, rw.Code)

	res := response{}

	err := json.Unmarshal(rw.Body.Bytes(), &res)
	assert.Nil(t, err)

	assert.Equal(t, "seed", res.Seed)
	assert.Equal(t, "address", res.Address)
}

// Test Get with error on provided seed
func TestGetErrorSeed(t *testing.T) {
	s := new(mockedService)
	router := createMockedService(s)
	req, rw := createRequest("GET", "/wallet?index=0", "")

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
}

// Test get with error in service
func TestGetErrorService(t *testing.T) {
	s := new(mockedService)
	s.On("Get", "seed", 0).Return(domain.Wallet{}, fmt.Errorf("An error ocurred while creating the wallet"))

	router := createMockedService(s)
	req, rw := createRequest("GET", "/wallet?seed=seed&index=0", "")

	router.ServeHTTP(rw, req)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
}
