package database

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/Lucas-Linhar3s/JobHub/pkg/config"
	"github.com/Lucas-Linhar3s/JobHub/pkg/log"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"
)

// Database is a struct that contains the database connection and the transaction timeout
type Database struct {
	tx                *sql.Tx
	db                *sql.DB
	transationTimeout int
	Builder           squirrel.StatementBuilderType
}

// NewDatabase is a function that returns a new database instance
func NewDatabase(config *config.Config, logger *log.Logger) *Database {
	db, err := open(config)
	if err != nil {
		logger.Fatal("failed to open database", zap.Error(err))
	}

	if err := db.db.Ping(); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}

	return db
}

func open(c *config.Config) (database *Database, err error) {
	var db *sql.DB

	switch c.Data.DB.User.Driver {
	case "sqlite":
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		path := filepath.Join(currentDir, "", c.Data.DB.User.Dsn)

		// Verifica se o arquivo do banco de dados existe
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, errors.New("banco de dados sqlite não encontrado")
		}

		db, err = sql.Open("sqlite3", path)
		if err != nil {
			return nil, err
		}
	case "postgres", "mysql":
		driverConfig := stdlib.DriverConfig{
			ConnConfig: pgx.ConnConfig{
				RuntimeParams: map[string]string{
					//Verificar
					"application_name": "github.com/Lucas-Linhar3s/Base-Structure-Golang",
					"DateStyle":        "ISO",
					"IntervalStyle":    "iso_8601",
					// TODO:
					"search_path": "public",
				},
			},
		}
		stdlib.RegisterDriverConfig(&driverConfig)

		db, err = sql.Open("pgx", driverConfig.ConnectionString(
			c.Data.DB.User.Nick+
				"://"+
				c.Data.DB.User.Username+
				":"+
				c.Data.DB.User.Password+
				"@"+
				c.Data.DB.User.HostName+
				":"+
				c.Data.DB.User.Port+
				"/"+
				c.Data.DB.User.Name))
		if err != nil {
			return nil, err
		}
	default:
		panic("unknown db driver")
	}

	db.SetMaxIdleConns(c.Data.DB.User.MaxIdle)
	db.SetMaxOpenConns(c.Data.DB.User.MaxConn)
	db.SetConnMaxLifetime(time.Second * 60)

	return &Database{
		db:                db,
		transationTimeout: c.Data.DB.User.TransationTimeout,
		Builder:           squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}, nil
}

// NewTransaction is a function that returns a new transaction instance
func (d *Database) NewTransaction() (*Database, error) {
	var (
		tx  *sql.Tx
		err error
	)

	// Define o timeout para a transação com `context.WithTimeout`
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.transationTimeout)*time.Second)
	defer cancel() // Certifique-se de cancelar o contexto quando sair da função


	tx, err = d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Database{
		tx:      tx,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(tx),
	}, nil
}

// Commit commits pending transactions for all databases open
func (t *Database) Commit() (err error) {
	err = t.tx.Commit()
	return
}

// Rollback is a function that rolls back the transaction
func (d *Database) Rollback() {
	_ = d.tx.Rollback()
}

// Close is a function that closes the database connection
func (d *Database) Close() error {
	if d.tx != nil {
		return errors.New("transaction not finished")
	}
	return d.db.Close()
}
