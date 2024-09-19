package main

import (
	"errors"
	"fmt"
)

func ConnectingToDb() error {
	fmt.Println("Connection wird erstellt")
	err := errors.New("fehler beim verbinden zur datenbank")

	return err
}

func main() {
	for i := 0; i < 2; i++ {
		defer fmt.Println("Connection wird geschlossen")

		err := ConnectingToDb()
		if err != nil {
			fmt.Println("Fehler:", err)
			continue
		}

		fmt.Println("Verbingung erfolgreich hergestellt")
	}
}
