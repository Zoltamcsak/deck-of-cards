package model

type CreateDeckRequest struct {
	Shuffled bool   `json:"shuffled"`
	Cards    string `json:"cards"`
}

type CreateDeckResponse struct {
	DeckId    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}
