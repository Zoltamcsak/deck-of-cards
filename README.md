# Deck of cards

## Rest API for cards operations

This project creates a REST API to handle the deck and cards to be used in 
games like Poker or Blackjack. 

It has the following 3 operations

### Create a new deck
`POST /decks`

It can accept query parameters of `cards` and `shuffled`

``
curl --request POST 'http://localhost:8080/decks?cards=AS,KD,AC,2C,KH&shuffled=true'
``

### Open a deck
`GET /decks/:id`

``
curl --request GET 'http://localhost:8080/decks/<deck-id>'
``

### Draw a card
`PUT /decks/:id/cards`

It can have a query parameter `count`

``
curl --request PUT 'http://localhost:8080/decks/<deck-id>/cards?count=3'
``

## Running the project

### Requirements
* Go 1.20 or above
* Docker

### Start the application locally
* Clone the project: `git clone https://github.com/Zoltamcsak/deck-of-cards.git`
* Run `make run-db` to start PostgreSQL locally (it'll connect to a DB called `deck_of_card` and uses port `5432`)
* Run `make tidy` to adjust dependencies
* Run `make run` to start the project locally
* You can access the application on port `:8080`

There's a `.env` file added to this repository just to make running locally easier. You can update any value there if needed!

There's a Postman collection called `deck_of_cards_postman_collection.json` that includes all 3 endpoints.

