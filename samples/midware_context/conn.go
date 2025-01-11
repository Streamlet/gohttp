package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	db *sql.DB
	tx *sql.Tx
}

func GetConnection(driverName string, dataSourceName string) (*Connection, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Connection{db: db}, nil
}

func (c *Connection) Close() {
	_ = c.db.Close()
}

func (c *Connection) Begin() error {
	if c.tx != nil {
		return errors.New("previous transaction not closed")
	}
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	c.tx = tx
	return nil
}

func (c *Connection) Commit() error {
	if c.tx == nil {
		return errors.New("not in transaction")
	}
	_ = c.tx.Commit()
	c.tx = nil
	return nil
}

func (c *Connection) Rollback() error {
	if c.tx == nil {
		return errors.New("not in transaction")
	}
	_ = c.tx.Rollback()
	c.tx = nil
	return nil
}
