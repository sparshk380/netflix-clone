package models

import (
	"context"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type ClientInterface interface {
	CloseConnection()
	InitializeDB() error
}

type Client struct {
	*pg.DB
}

func NewDBClient() Client {
	// tlsConfig := tls.Config{
	// 	InsecureSkipVerify: true,
	// }

	db := pg.Connect(&pg.Options{
		Addr:     os.Getenv("DB_ADDR"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB"),
		// TLSConfig: &tlsConfig,
		// PoolSize:
	})

	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}

	return Client{
		db,
	}
}

func (db *Client) CloseConnection() {
	db.Close()
}

func (db *Client) InitializeDB() error {
	logrus.Infoln("Starting DB initialization")

	if err := db.CreateAccountsSchema(); err != nil {
		return err
	}

	if err := db.CreateNodesSchema(); err != nil {
		return err
	}

	return nil
}
