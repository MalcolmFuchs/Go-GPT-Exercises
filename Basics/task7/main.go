package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type DatabaseConnection struct {
	Name string
}

func (db DatabaseConnection) Connect() error {
	fmt.Printf("Verbindung zu %s wird aufgebaut...\n", db.Name)
	time.Sleep(time.Millisecond * 500)

	if db.Name == "MySQL" {
		return errors.New("Fehler bei der Verbindung zur Datenbank: " + db.Name)
	}

	fmt.Printf("Verbindung zu %s erfolgreich aufgebaut\n", db.Name)
	return nil
}

func handleDatabaseConnection(db DatabaseConnection, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	if err := db.Connect(); err != nil {
		errChan <- err
		return
	}
	errChan <- nil
}

func main() {
	databases := []DatabaseConnection{
		{Name: "CockroachDB"},
		{Name: "MySQL"},
		{Name: "PostgreSQL"},
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(databases))

	for _, db := range databases {
		wg.Add(1)
		go handleDatabaseConnection(db, &wg, errChan)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			fmt.Println("Fehler:", err)
		}
	}

	fmt.Println("Alle Verbindungen zu den Datenbanken wurden abgeschlossen.")
}
