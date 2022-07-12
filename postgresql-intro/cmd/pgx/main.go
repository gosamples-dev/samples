package main

import (
	"context"
	"log"
	"postgresql-intro/app"
	"postgresql-intro/website"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	dbpool, err := pgxpool.Connect(context.Background(), "postgres://postgres:mysecretpassword@localhost:5432/website")
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	websiteRepository := website.NewPostgreSQLPGXRepository(dbpool)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.RunRepositoryDemo(ctx, websiteRepository)
}
