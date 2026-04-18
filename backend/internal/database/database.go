package database

import (
	"github.com/dhruv15803/budgeting-app/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	dbConfig config.DbConfig
}

func NewPostgres(dbConfig config.DbConfig) *Postgres {
	return &Postgres{
		dbConfig: dbConfig,
	}
}

func (p *Postgres) Connect() (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", p.dbConfig.Url)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(p.dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(p.dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(p.dbConfig.MaxConnLifetime)
	db.SetConnMaxIdleTime(p.dbConfig.MaxConnIdleTime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
