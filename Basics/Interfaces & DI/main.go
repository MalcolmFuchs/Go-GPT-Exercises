package main

import (
	"fmt"
	"log"
	"reflect"
)

type DataStore interface {
	GetData() (string, error)
}

type DataReader struct{}

func (d DataReader) GetData() (string, error) {
	return "Daten werden besorgt", nil
}

type InMemoryDataReader struct{}

func (m InMemoryDataReader) GetData() (string, error) {
	return "Daten aus der In-Memory-Quelle werden besorgt", nil
}

func handleGetData(ds DataStore) {
	log.Println("Verwende Datenquelle:", reflect.TypeOf(ds).String())

	message, err := ds.GetData()
	if err != nil {
		log.Println("Fehler beim Abrufen der Daten:", err)
		return
	}
	fmt.Println(message)
}

func main() {
	inMemory := &InMemoryDataReader{}
	dataReader := &DataReader{}

	handleGetData(inMemory)
	handleGetData(dataReader)
}
