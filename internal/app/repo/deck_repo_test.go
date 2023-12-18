package repo

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

// TestContainer represents a test container for the PostgreSQL database.
type TestContainer struct {
	container testcontainers.Container
}

// NewTestContainer creates and starts a new PostgreSQL test container.
func NewTestContainer() (*TestContainer, error) {
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		postgres.WithDatabase("deck_of_card"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	return &TestContainer{container: container}, nil
}

// Close stops and removes the test container.
func (tc *TestContainer) Close() error {
	return tc.container.Terminate(context.Background())
}

// GetDSN returns the data source name for connecting to the test container's PostgreSQL database.
func (tc *TestContainer) GetDSN() string {
	host, err := tc.container.Host(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	port, err := tc.container.MappedPort(context.Background(), "5432")
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("user=postgres password=postgres dbname=postgres sslmode=disable host=%s port=%s", host, port.Port())
}

func setupTestContainer(t *testing.T) (*sqlx.DB, *TestContainer, func()) {
	container, err := NewTestContainer()
	assert.NoError(t, err)

	dsn := container.GetDSN()
	db, err := sqlx.Open("postgres", dsn)
	assert.NoError(t, err)

	// should be created from an init sql script or migrate script
	db.Exec(`
		create table if not exists decks (
			id varchar(50) primary key,
			shuffled bool default false not null,
			remaining int not null,
			cards text[] not null,
			created_at timestamp without time zone default current_timestamp not null,
			updated_at timestamp without time zone default current_timestamp not null
		)
	`)

	return db, container, func() {
		assert.NoError(t, db.Close())
		assert.NoError(t, container.Close())
	}
}

func TestDeckRepo(t *testing.T) {
	db, _, cleanup := setupTestContainer(t)
	defer cleanup()

	repo := NewDeckRepo(db)

	// Test case: CreateDeck
	deck := Deck{
		Id:        "test-deck-id",
		Shuffled:  true,
		Remaining: 52,
		Cards:     []string{"AH", "2C", "3D", "4S", "5H"},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := repo.CreateDeck(deck)
	assert.NoError(t, err)

	// Test case: GetDeckById
	fetchedDeck, err := repo.GetDeckById("test-deck-id")
	assert.NoError(t, err)
	assert.NotNil(t, fetchedDeck)
	assert.Equal(t, deck.Id, fetchedDeck.Id)

	// Test case: UpdateDeck
	deck.Remaining = 50
	err = repo.UpdateDeck(deck)
	assert.NoError(t, err)

	// Verify the updated deck
	updatedDeck, err := repo.GetDeckById("test-deck-id")
	assert.NoError(t, err)
	assert.NotNil(t, updatedDeck)
	assert.Equal(t, deck.Remaining, updatedDeck.Remaining)
}
