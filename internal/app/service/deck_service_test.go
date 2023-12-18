package service

import (
	"database/sql"
	"errors"
	"fmt"
	customErr "github.com/deck/internal/app/error"
	"github.com/deck/internal/app/model"
	"github.com/deck/internal/app/repo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sort"
	"testing"
	"time"
)

var sequentialDeck = []string{
	"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "JS", "QS", "KS",
	"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "JD", "QD", "KD",
	"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "10C", "JC", "QC", "KC",
	"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "JH", "QH", "KH",
}

// MockRepo is a mock implementation of the DeckRepo interface
type MockRepo struct {
	Decks     map[string]repo.Deck
	DeckError error
}

// Implement the DeckRepo interface methods for the mock
func (m *MockRepo) CreateDeck(deck repo.Deck) error {
	m.Decks[deck.Id] = deck
	return nil
}

func (m *MockRepo) GetDeckById(id string) (*repo.Deck, error) {
	deck, found := m.Decks[id]
	if m.DeckError != nil {
		return nil, m.DeckError
	}
	if !found {
		return nil, sql.ErrNoRows
	}
	return &deck, nil
}

func (m *MockRepo) UpdateDeck(deck repo.Deck) error {
	m.Decks[deck.Id] = deck
	return m.DeckError
}

func TestCreateDeck(t *testing.T) {
	mockRepo := &MockRepo{Decks: make(map[string]repo.Deck)}
	deckService := NewDeckService(mockRepo)

	// Test case: default deck creation
	req := model.CreateDeckRequest{Shuffled: false}
	res, err := deckService.CreateDeck(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 52, res.Remaining)

	// Test case: custom cards
	req = model.CreateDeckRequest{Cards: "AS,2S,3S", Shuffled: false}
	res, err = deckService.CreateDeck(req)

	// Test case: default deck with shuffled cards
	req = model.CreateDeckRequest{Shuffled: true}
	res, err = deckService.CreateDeck(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 52, res.Remaining)
	assert.Equal(t, true, res.Shuffled)

	// Test case: create a deck with an invalid card code
	req = model.CreateDeckRequest{
		Cards: "XYZ",
	}
	res, err = deckService.CreateDeck(req)
	assert.Error(t, err)
	assert.Nil(t, res)

	// Test case: create a deck with duplicate cards
	req = model.CreateDeckRequest{
		Cards: "2H,2H,3C,4D",
	}
	res, err = deckService.CreateDeck(req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestGetDeckById(t *testing.T) {
	// Set up the DeckService with the mock repository
	mockRepo := &MockRepo{Decks: make(map[string]repo.Deck)}
	deckService := NewDeckService(mockRepo)
	now := time.Now().UTC()

	// Test case: get a deck by ID successfully
	deckID := "existing_deck_id"
	mockDeck := repo.Deck{
		Id:        deckID,
		Shuffled:  true,
		Remaining: 3,
		Cards:     []string{"AH", "2C", "3D"},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockRepo.Decks[deckID] = mockDeck

	res, err := deckService.GetDeckById(deckID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, deckID, res.DeckId)
	assert.Equal(t, true, res.Shuffled)
	assert.Equal(t, 3, res.Remaining)
	assert.Len(t, res.Cards, 3)

	// Test case: get a non-existing deck by ID
	nonExistingDeckID := "non_existing_deck_id"
	res, err = deckService.GetDeckById(nonExistingDeckID)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("deck with id %s wasn't found", nonExistingDeckID))
	assert.Nil(t, res)

	// Test case: error while retrieving deck from the database
	errMessage := "database error"
	mockRepo.DeckError = errors.New(errMessage)
	res, err = deckService.GetDeckById(deckID)

	assert.Error(t, err)
	assert.EqualError(t, err, errMessage)
	assert.Nil(t, res)
}

func TestDrawCards(t *testing.T) {
	// Set up the DeckService with the mock repository
	mockRepo := &MockRepo{Decks: make(map[string]repo.Deck)}
	deckService := NewDeckService(mockRepo)
	now := time.Now().UTC()

	// Test case: draw cards successfully
	deckID := "existing_deck_id"
	initialRemaining := 5
	mockDeck := repo.Deck{
		Id:        deckID,
		Shuffled:  true,
		Remaining: initialRemaining,
		Cards:     []string{"AH", "2C", "3D", "4S", "5H"},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockRepo.Decks[deckID] = mockDeck

	count := 3
	cards, err := deckService.DrawCards(deckID, count)

	assert.NoError(t, err)
	assert.Len(t, cards, count)
	updatedDeck, found := mockRepo.Decks[deckID]
	assert.True(t, found)
	assert.Equal(t, initialRemaining-count, updatedDeck.Remaining)

	// Test case: draw cards with count exceeding remaining
	count = 15
	cards, err = deckService.DrawCards(deckID, count)

	assert.Error(t, err)
	assert.EqualError(t, err, "count must be less or equal than deck's remaining")
	assert.Nil(t, cards)

	// Test case: draw cards with invalid count
	count = 0
	cards, err = deckService.DrawCards(deckID, count)

	assert.Error(t, err)
	assert.EqualError(t, err, "count must be between 1 - 52")
	assert.Nil(t, cards)

	// Test case: error while updating the deck
	errMessage := "update error"
	mockRepo.DeckError = errors.New(errMessage)
	count = 2
	cards, err = deckService.DrawCards(deckID, count)

	assert.Error(t, err)
	assert.Nil(t, cards)
}

func TestGenerateDefaultDeck(t *testing.T) {
	// Test case: generate the default deck
	result := GenerateDefaultDeck()

	assert.Equal(t, sequentialDeck, result)
}

func TestShuffleCards(t *testing.T) {
	// Test case: shuffle a deck of cards
	// Make a copy of the original deck to compare
	expectedDeck := make([]string, len(sequentialDeck))
	copy(expectedDeck, sequentialDeck)

	ShuffleCards(sequentialDeck)

	// Assert that the deck is shuffled by comparing it with the original deck
	assert.NotEqual(t, expectedDeck, sequentialDeck)

	// Sort both decks to compare the elements
	sort.Strings(expectedDeck)
	sort.Strings(sequentialDeck)

	// Assert that the elements are the same after sorting, indicating a valid shuffle
	assert.Equal(t, expectedDeck, sequentialDeck)
}

func TestValidateCards(t *testing.T) {
	// Test case: valid cards
	validCards := []string{"AH", "2C", "3D", "4S", "5H"}

	err := validateCards(validCards)

	assert.NoError(t, err)

	// Test case: invalid card code
	invalidCard := "INVALID"
	invalidCards := []string{"AH", invalidCard, "3D", "4S", "5H"}

	err = validateCards(invalidCards)

	assert.Error(t, err)
	assert.EqualError(t, err, customErr.New(http.StatusBadRequest, "contains invalid card code").Error())

	// Test case: duplicate cards
	duplicateCard := "3D"
	duplicateCards := []string{"AH", "2C", duplicateCard, duplicateCard, "5H"}

	err = validateCards(duplicateCards)

	assert.Error(t, err)
	assert.EqualError(t, err, customErr.New(http.StatusBadRequest, "contains duplicate").Error())
}

// Additional test for isValidCardCode function
func TestIsValidCardCode(t *testing.T) {
	// Test case: valid card code
	validCard := "2H"

	result := isValidCardCode(validCard)

	assert.True(t, result)

	// Test case: invalid rank in card code
	invalidRankCard := "XH"

	result = isValidCardCode(invalidRankCard)

	assert.False(t, result)

	// Test case: invalid suit in card code
	invalidSuitCard := "2X"

	result = isValidCardCode(invalidSuitCard)

	assert.False(t, result)
}

func TestGetValueAndSuit(t *testing.T) {
	// Test case: valid card code
	validCard := "2H"

	card, err := getValueAndSuit(validCard)

	// Assert that there is no error
	assert.NoError(t, err)

	expectedCard := &model.Card{
		Value: repo.Values[("2")],
		Suit:  repo.Suites[("H")],
		Code:  validCard,
	}
	assert.Equal(t, expectedCard, card)
}
