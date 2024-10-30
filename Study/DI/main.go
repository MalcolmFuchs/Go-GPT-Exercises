// Clean Architecture & Dependency Injection

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type DataReader interface {
	ReadData() ([]string, error)
}

type FileDataReader struct {
	FileName string
}

func (f FileDataReader) ReadData() ([]string, error) {
	file, err := os.Open(f.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	log.Println("Datei erfolgreich geöffnet")

	var data []string
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		data = append(data, scan.Text())
	}

	return data, nil
}

type InMemoryDataReader struct {
	Data []string
}

func (m InMemoryDataReader) ReadData() ([]string, error) {
	return m.Data, nil
}

func ProcessData(reader DataReader) {
	data, err := reader.ReadData()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Gelesene Daten:")

	for _, d := range data {
		fmt.Println(d)
	}
}

func main() {

	fileReader := &FileDataReader{FileName: "poem.txt"}
	ProcessData(fileReader)
}
