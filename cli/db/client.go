package db

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Client is a struct that holds the connection pool
type Client struct {
	connPool *pgxpool.Pool
}

// NewClient creates a new client and returns it
func NewClient() *Client {
	return &Client{
		connPool: nil,
	}
}

// Close closes the connection pool
func (c *Client) Close() {
	c.connPool.Close()
}

// Connect connects to the database and returns a client
func (c *Client) Connect(ctx context.Context, connStr string) error {
	if c.connPool != nil {
		log.Debug("Connection pool already exists, returning existing client")
		return nil
	}

	connPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return err
	}

	c.connPool = connPool

	return nil
}
