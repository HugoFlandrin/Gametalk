package main

import (
	api "api/api"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func AddCategories() {
	db := api.OpenDatabase("BDD.db")
	defer db.Close()

	categories := []struct {
		Name        string
		Description string
	}{
		{"Wiki", "Encyclopedia of information"},
		{"Tournois", "Details of tournaments"},
		{"A Venir", "Upcoming events"},
		{"Rumeur", "Rumors and speculations"},
		{"Jeux", "Information about games"},
		{"Message", "Messages and communications"},
	}

	stmt, err := db.Prepare("INSERT INTO CATEGORIES(name, description) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, category := range categories {
		_, err := stmt.Exec(category.Name, category.Description)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Categories added successfully!")
}

func ClearLike() {
	db := api.OpenDatabase("BDD.db")
	defer db.Close()
	db.Exec("DELETE FROM DISLIKES WHERE id > 0")
}

func AddTopicIDColumn() {
	db := api.OpenDatabase("BDD.db")
	defer db.Close()

	alterTableSQL := `ALTER TABLE DISLIKES ADD COLUMN topic_id INTEGER;`

	_, err := db.Exec(alterTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Column topic_id added successfully!")
}

func main() {
	ClearLike()
}
