package database

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// From psqldb/docker-compose.yml and psqldb/init.sql
const (
	POSTGRES_URL = "postgresql://jirauser:asdf@localhost:5432/jiradb?sslmode=disable"
)

type PsqlInterface struct {
	db *sqlx.DB
}

func NewPsqlInterface() *PsqlInterface {
	p := &PsqlInterface{}
	db, err := sqlx.Connect("postgres", POSTGRES_URL)

	if err != nil {
		log.Panicln(err)
	}

	p.db = db

	return p
}

func (p *PsqlInterface) GetItem(name string) (i *Item) {
	i = &Item{}
	err := p.db.Get(i, "SELECT * FROM Items WHERE Name=$1", name)

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		log.Panicln(err)
	}

	return
}
