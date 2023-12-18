package service

import (
	"database/sql"
	"fmt"
	customErr "github.com/deck/internal/app/error"
	"github.com/deck/internal/app/model"
	"github.com/deck/internal/app/repo"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type DeckService struct {
	repo *repo.DeckRepo
}

func NewDeckService(repo *repo.DeckRepo) *DeckService {
	return &DeckService{repo: repo}
}

func (s *DeckService) CreateDeck(req model.CreateDeckRequest) (*model.CreateDeckResponse, error) {
	var cards []string
	if len(req.Cards) == 0 {
		cards = GenerateDefaultDeck()
	} else {
		cards = strings.Split(req.Cards, ",")
		err := validateCards(cards)
		if err != nil {
			return nil, err
		}
	}
	if req.Shuffled {
		ShuffleCards(cards)
	}
	now := time.Now().UTC()
	deck := repo.Deck{
		Id:        uuid.New().String(),
		Shuffled:  req.Shuffled,
		Remaining: len(cards),
		Cards:     cards,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := s.repo.CreateDeck(deck)
	if err != nil {
		return nil, customErr.Wrap(http.StatusInternalServerError, "couldn't save deck", err)
	}
	return &model.CreateDeckResponse{
		DeckId:    deck.Id,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
	}, nil
}

func (s *DeckService) GetDeckById(id string) (*model.OpenDeckResponse, error) {
	deck, err := s.repo.GetDeckById(id)
	if err == sql.ErrNoRows {
		return nil, customErr.New(http.StatusNotFound, fmt.Sprintf("deck with id %s wasn't found", id))
	}
	if err != nil {
		return nil, customErr.Wrap(http.StatusInternalServerError, "couldn't get deck from the database", err)
	}
	cards := make([]model.Card, len(deck.Cards))
	for i, c := range deck.Cards {
		card, err := getValueAndSuit(c)
		if err != nil {
			return nil, err
		}
		cards[i] = *card
	}
	return &model.OpenDeckResponse{
		DeckId:    deck.Id,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
		Cards:     cards,
	}, nil
}

func (s *DeckService) DrawCards(id string, count int) ([]model.Card, error) {
	if count <= 0 || count > 52 {
		return nil, customErr.New(http.StatusBadRequest, "count must be between 1 - 52")
	}
	deck, err := s.GetDeckById(id)
	if err != nil {
		return nil, err
	}
	if count > deck.Remaining {
		return nil, customErr.New(http.StatusBadRequest, "count must be less or equal than deck's remaining")
	}
	cards := drawFirstCards(*deck, count)
	updatedDeck := updateDeck(*deck, count)
	err = s.repo.UpdateDeck(updatedDeck)
	if err != nil {
		return nil, customErr.Wrap(http.StatusInternalServerError, "couldn't update deck", err)
	}
	return cards, nil
}

func GenerateDefaultDeck() []string {
	var deck []string
	for _, s := range repo.SequentialSuits {
		for _, v := range repo.SequentialValues {
			deck = append(deck, fmt.Sprintf("%s%s", v, s))
		}
	}
	return deck
}

func ShuffleCards(cards []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Seed(time.Now().UnixNano())
	r.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

func updateDeck(deck model.OpenDeckResponse, count int) repo.Deck {
	cardCodes := make([]string, len(deck.Cards)-count)
	for i, c := range deck.Cards[count:] {
		cardCodes[i] = c.Code
	}
	updatedDeck := repo.Deck{
		Id:        deck.DeckId,
		Remaining: deck.Remaining - count,
		Shuffled:  deck.Shuffled,
		Cards:     cardCodes,
	}
	return updatedDeck
}

func validateCards(cards []string) error {
	checkDuplicates := make(map[string]bool, len(cards))
	for _, c := range cards {
		if !isValidCardCode(c) {
			return customErr.New(http.StatusBadRequest, "contains invalid card code")
		}
		if _, found := checkDuplicates[c]; found {
			return customErr.New(http.StatusBadRequest, "contains duplicate")
		}
		checkDuplicates[c] = true
	}
	return nil
}

func isValidCardCode(code string) bool {
	validRanks := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		repo.Ace, repo.Two, repo.Three, repo.Four, repo.Five, repo.Six, repo.Seven, repo.Eight, repo.Nine, repo.Ten, repo.Jack, repo.Queen, repo.King)
	validSuits := fmt.Sprintf("%s|%s|%s|%s", repo.Spades, repo.Hearts, repo.Diamonds, repo.Clubs)

	validCardPattern := regexp.MustCompile(fmt.Sprintf("^(%s)(%s)$", validRanks, validSuits))
	return validCardPattern.MatchString(strings.ToUpper(code))
}

func getValueAndSuit(code string) (*model.Card, error) {
	if !isValidCardCode(code) {
		return nil, customErr.New(http.StatusInternalServerError, "code is not valid")
	}
	valueCode := code[:len(code)-1]
	suitCode := string(code[len(code)-1])

	fullValue := repo.Values[repo.CardCode(valueCode)]
	suit := repo.Suites[repo.SuitCode(suitCode)]

	return &model.Card{
		Value: fullValue,
		Suit:  suit,
		Code:  code,
	}, nil
}

func drawFirstCards(deck model.OpenDeckResponse, count int) []model.Card {
	return deck.Cards[:count]
}
