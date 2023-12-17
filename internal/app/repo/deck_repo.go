package repo

import (
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"time"
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

func (r *DeckRepo) UpdateDeck(deck Deck) error {
	res, err := r.db.Exec(`update decks set shuffled=$1, remaining=$2, cards=$3, updated_at=$4 where id=$5`,
		deck.Shuffled, deck.Remaining, deck.Cards, time.Now().UTC(), deck.Id)
	if err != nil {
		glog.Errorf("error while updating deck with id %s", deck.Id, err)
		return err
	}
	rows, err := res.RowsAffected()
	glog.Infof("%d rows updated", rows)
	return nil
}
