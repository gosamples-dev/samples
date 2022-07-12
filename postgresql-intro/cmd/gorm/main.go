package main

import (
	"context"
	"log"
	"postgresql-intro/app"
	"postgresql-intro/website"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gormDB, err := gorm.Open(postgres.Open("postgres://postgres:mysecretpassword@localhost:5432/website"))
	if err != nil {
		log.Fatal(err)
	}

	websiteRepository := website.NewPostgreSQLGORMRepository(gormDB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.RunRepositoryDemo(ctx, websiteRepository)
}
