package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gosamples-dev/samples/sqlite-intro/website"
	_ "github.com/mattn/go-sqlite3"
)

const fileName = "sqlite.db"

func main() {
	os.Remove(fileName)

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}

	websiteRepository := website.NewSQLiteRepository(db)

	if err := websiteRepository.Migrate(); err != nil {
		log.Fatal(err)
	}

	gosamples := website.Website{
		Name: "GOSAMPLES",
		URL:  "https://gosamples.dev",
		Rank: 2,
	}
	golang := website.Website{
		Name: "Golang official website",
		URL:  "https://golang.org",
		Rank: 1,
	}

	createdGosamples, err := websiteRepository.Create(gosamples)
	if err != nil {
		log.Fatal(err)
	}
	createdGolang, err := websiteRepository.Create(golang)
	if err != nil {
		log.Fatal(err)
	}

	gotGosamples, err := websiteRepository.GetByName("GOSAMPLES")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("get by name: %+v\n", gotGosamples)

	createdGosamples.Rank = 1
	if _, err := websiteRepository.Update(createdGosamples.ID, *createdGosamples); err != nil {
		log.Fatal(err)
	}

	all, err := websiteRepository.All()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nAll websites:\n")
	for _, website := range all {
		fmt.Printf("website: %+v\n", website)
	}

	if err := websiteRepository.Delete(createdGolang.ID); err != nil {
		log.Fatal(err)
	}

	all, err = websiteRepository.All()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nAll websites:\n")
	for _, website := range all {
		fmt.Printf("website: %+v\n", website)
	}
}
