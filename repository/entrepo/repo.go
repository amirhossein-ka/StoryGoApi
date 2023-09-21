package entrepo

import (
	"StoryGoAPI/config"
	"StoryGoAPI/ent"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repo struct {
	client *ent.Client
}

func (r *Repo) Close() error {
	return r.client.Close()
}

// DropAllTables as its name says drops all tables, DO NOT RUN IT ANYWHERE. IT IS FOR TESTS ONLY.
func (r *Repo) DropAllTables() error {
	if _, err := r.client.GuestUser.Delete().Exec(context.Background()); err != nil {
		return err
	}
	if _, err := r.client.User.Delete().Exec(context.Background()); err != nil {
		return err
	}
	if _, err := r.client.Story.Delete().Exec(context.Background()); err != nil {
		return err
	}
	return nil
}

func NewRepo(cfg config.DataBase) (*Repo, error) {
	switch cfg.Driver {
	case "pgx":
		cli, err := openPgx(cfg.URI)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to db: %w", err)
		}

		if err = cli.Schema.Create(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to create schema: %w", err)
		}

		return &Repo{client: cli}, nil
	case "mysql":

		cli, err := retryConnect(cfg.URI, 5, time.Second*5)
		if err != nil {
			return nil, err
		}
		if err = cli.Schema.Create(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to create schema: %w", err)
		}
		return &Repo{client: cli}, nil

	default:
		return nil, fmt.Errorf("invalid database driver")
	}
}

// since I could not figure it how to do a healthcheck in docker-compose,
// I write the healthcheck logic here,
// it retries 3 times with 5-second delays and a delay for start
func retryConnect(uri string, maxRetries int, retryInterval time.Duration) (*ent.Client, error) {
	var (
		cli *ent.Client
		err error
	)
	for i := 0; i < maxRetries; i++ {
		cli, err = openMysql(uri)
		if err == nil {
			break
		}
		log.Printf("error connecting database: %v, retrying...\n", err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, err
	}

	return cli, nil
}

func openPgx(uri string) (*ent.Client, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}

func openMysql(uri string) (*ent.Client, error) {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	drv := entsql.OpenDB(dialect.MySQL, db)
	return ent.NewClient(ent.Driver(drv)), nil
}
