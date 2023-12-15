package repo

import (
	"github.com/lib/pq"
	"time"
)

type CardCode string
type SuitCode string

const (
	Ace   CardCode = "A"
	Two   CardCode = "2"
	Three CardCode = "3"
	Four  CardCode = "4"
	Five  CardCode = "5"
	Six   CardCode = "6"
	Seven CardCode = "7"
	Eight CardCode = "8"
	Nine  CardCode = "9"
	Ten   CardCode = "10"
	Jack  CardCode = "J"
	Queen CardCode = "Q"
	King  CardCode = "K"
)

const (
	Spades   SuitCode = "S"
	Diamonds SuitCode = "D"
	Clubs    SuitCode = "C"
	Hearts   SuitCode = "H"
)

var SequentialValues = []CardCode{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}
var SequentialSuits = []SuitCode{Spades, Diamonds, Clubs, Hearts}

type Deck struct {
	Id        string         `db:"id"`
	Shuffled  bool           `db:"shuffled"`
	Remaining int            `db:"remaining"`
	Cards     pq.StringArray `db:"cards"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
