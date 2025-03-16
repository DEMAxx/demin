package migrations

import (
	"context"
	"embed"
	"fmt"
	"github.com/DEMAxx/demin/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pressly/goose/v3"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

//go:embed sources/*.sql
var embedMigrations embed.FS

func Run(cnf *Config) error {
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.Name)

	st, err := sqlstorage.New(dsn)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := st.Connect(ctx); err != nil {
		return err
	}

	defer func(st *sqlstorage.Storage, ctx context.Context) {
		err := st.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(st, ctx)

	return goose.Up(st.GetDb(), "sources")
}
