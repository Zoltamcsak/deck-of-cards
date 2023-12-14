package config

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"os"
)

type Database struct {
	host          string
	user          string
	pass          string
	name          string
	migrationPath string
}

func NewDbConnection() (*sqlx.DB, error) {
	dbConf := Database{
		host:          os.Getenv("DB_HOST"),
		user:          os.Getenv("DB_USER"),
		pass:          os.Getenv("DB_PASSWORD"),
		name:          os.Getenv("DB_NAME"),
		migrationPath: os.Getenv("MIGRATION_PATH"),
	}

	db, err := sqlx.Connect("postgres", dbConf.sourceName())
	if err != nil {
		return nil, err
	}
	if err = dbConf.doMigrations(db.DB); err != nil {
		return nil, err
	}

	return db, nil
}

func (db Database) sourceName() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", db.host, db.user, db.pass, db.name)
}

func (db Database) doMigrations(sqlDb *sql.DB) error {
	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(db.migrationPath, db.name, driver)

	if err != nil {
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		glog.Infof("no migration required")
		return nil
	}

	return err
}
