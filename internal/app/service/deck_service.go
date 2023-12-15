package service

import (
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
		cards = generateDefaultDeck()
	} else {
		cards = strings.Split(req.Cards, ",")
		err := validateCards(cards)
		if err != nil {
			return nil, err
		}
	}
	if req.Shuffled {
		shuffleCards(cards)
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

func generateDefaultDeck() []string {
	var deck []string
	for _, s := range repo.SequentialSuits {
		for _, v := range repo.SequentialValues {
			deck = append(deck, fmt.Sprintf("%s%s", v, s))
		}
	}
	return deck
}

func shuffleCards(cards []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Seed(time.Now().UnixNano())
	r.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
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