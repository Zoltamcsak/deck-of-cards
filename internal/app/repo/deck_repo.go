package repo

import (
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
)

type DeckRepo struct {
	db *sqlx.DB
}

func NewCardRepo(db *sqlx.DB) *DeckRepo {
	return &DeckRepo{db: db}
}

func (r *DeckRepo) CreateDeck(deck Deck) error {
	_, err := r.db.NamedExec(`insert into decks (id, shuffled, remaining, cards, created_at, updated_at) 
                          values (:id, :shuffled, :remaining, :cards, :created_at, :updated_at)`, deck)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeckRepo) GetDeckById(id string) (*Deck, error) {
	var deck Deck
	err := r.db.Get(&deck, "select * from decks where id=$1", id)
	if err != nil {
		glog.Errorf("error while getting deck with id %s", id, err)
		return nil, err
	}
	return &deck, nil
}
