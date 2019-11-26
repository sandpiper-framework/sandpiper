package postgres

import (
	"context"
	"log"
	"time"

	"github.com/go-pg/pg/v9"
	// DB adapter
	_ "github.com/lib/pq"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(ctx context.Context, pg *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	log.Printf(q.FormattedQuery())
	return nil
}

// New creates new database connection to a postgres database
func New(psn string, timeout int, enableLog bool) (*pg.DB, error) {
	u, err := pg.ParseURL(psn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(u)

	_, err = db.Exec("SELECT 1")	// test connectivity
	if err != nil {
		return nil, err
	}

	if timeout > 0 {
		db.WithTimeout(time.Second * time.Duration(timeout))
	}

	if enableLog {
		db.AddQueryHook(dbLogger{})
	}

	return db, nil
}
