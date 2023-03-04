package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Client communicates with the database.
type Client struct {
	url  string
	ctx  context.Context
	conn *pgx.Conn
}

// Connect to the database and return the connection
func NewDataBase(ctx context.Context, url string) *Client {
	var db Client
	var err error

	db.url = url
	db.ctx = ctx
	db.conn, err = pgx.Connect(db.ctx, db.url)
	if err != nil {
		panic(err)
	}

	return &db
}

// DBQueryArgs
type QueryArg struct {
	Sql  string
	Args []interface{}
}

// Do a query to the database
func (db *Client) Query(query QueryArg, v interface{}) ([]interface{}, error) {
	// Query DB
	rows, err := db.conn.Query(db.ctx, query.Sql, query.Args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the result
	var result []interface{}
	for rows.Next() {
		err = rows.Scan(v)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
}

// QueryRow and return the result
func (db *Client) Exec(query QueryArg, dest ...interface{}) error {
	// Query and return the result
	err := db.conn.QueryRow(db.ctx, query.Sql, query.Args...).Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}
