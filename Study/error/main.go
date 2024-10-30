// Error-Handling

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {

	file, err := os.Open("poem.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.Println("Datei erfolgreich geöffnet")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
