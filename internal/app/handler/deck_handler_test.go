package handler

import (
	"encoding/json"
	custErr "github.com/deck/internal/app/error"
	"github.com/deck/internal/app/model"
	"github.com/deck/internal/app/repo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	Decks     map[string]repo.Deck
	DeckError error
}

func (m *MockService) CreateDeck(req model.CreateDeckRequest) (*model.CreateDeckResponse, error) {
	if m.DeckError != nil {
		return nil, m.DeckError
	}
	return nil, nil
}
func (m *MockService) GetDeckById(id string) (*model.OpenDeckResponse, error) {
	if m.DeckError != nil {
		return nil, m.DeckError
	}
	return &model.OpenDeckResponse{
		DeckId:    "valid-deck-id",
		Shuffled:  true,
		Remaining: 1,
		Cards: []model.Card{{
			Value: "ACE",
			Suit:  "CLUBS",
			Code:  "AC",
		}},
	}, nil
}
func (m *MockService) DrawCards(id string, count int) ([]model.Card, error) {
	return []model.Card{{Value: "A", Suit: "Spades", Code: "AS"}}, nil
}

var router *gin.Engine
var mockService *MockService

func init() {
	gin.SetMode(gin.TestMode)

	// Create a mock service with a mock repository
	mockService = &MockService{}
	deckHandler := NewDeckHandler(mockService)

	// Create a test Gin router and add the CreateDeck route
	router = gin.Default()
	deckHandler.InitRoutes(router)
}

func TestCreateDeckHandler(t *testing.T) {
	// Test case: Create a deck with default parameters
	w := performRequest(router, "POST", "/decks", "")
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test case: Create a deck with cards parameter
	w = performRequest(router, "POST", "/decks?cards=AS,KD,AC,2C,KH", "")
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test case: Create a deck with shuffled parameter
	w = performRequest(router, "POST", "/decks?shuffled=true", "")
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test case: Create a deck with invalid shuffled parameter
	w = performRequest(router, "POST", "/decks?shuffled=invalid", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDeckByIdHandler(t *testing.T) {
	// Test case: Get a deck by valid ID
	w := performRequest(router, "GET", "/decks/valid-deck-id", "")
	var actualResult model.OpenDeckResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResult)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "valid-deck-id", actualResult.DeckId)

	// Test case: Get a deck by invalid ID
	mockService.DeckError = custErr.New(http.StatusNotFound, "not found")
	w = performRequest(router, "GET", "/decks/invalid-deck-id", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDrawCardsHandler(t *testing.T) {
	// Test case: Draw cards with a valid count
	w := performRequest(router, "PUT", "/decks/valid-deck-id/cards?count=3", "")
	assert.Equal(t, http.StatusCreated, w.Code)

	// Test case 2: Draw cards with an invalid count
	w = performRequest(router, "PUT", "/decks/valid-deck-id/cards?count=invalid", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// performRequest is a helper function to send a request to the Gin router and return the response recorder.
func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
