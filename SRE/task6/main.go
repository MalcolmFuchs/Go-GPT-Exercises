package main

import (
	"errors"
	"fmt"
)

var MockDatabase = []string{"cockroach", "postgres", "mySQL"}

func OpenDatabase(db string) error {
	fmt.Printf("Verbindung zu %s wurde aufgebaut \n", db)

	if db == "mySQL" {
		return errors.New("Fehler beim Vebindung zur Datenbank:")
	}

	return nil
}

func main() {

	for _, db := range MockDatabase {

		func(db string) {
			err := OpenDatabase(db)
			if err != nil {
				fmt.Println("Fehler:", err)
				return
			}

			defer fmt.Println("Verbindung zur Datenbank wird geschlossen")
			fmt.Println("Verbindung wurde erfolgreich aufgebaut")
		}(db)

	}
}
