package repo

import "github.com/jmoiron/sqlx"

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
