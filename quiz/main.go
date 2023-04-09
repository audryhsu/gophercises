package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// open csv file
	// todo: ask for a file path of csv or else read in problems.csv as default
	var csvFile string = "problems.csv"
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	// NewReader takes an io.Reader, and file implements the io.Reader interface
	csvReader := csv.NewReader(file)
	csvReader.Comma = ','
	for {
		// loop over each record until EOF
		// a record is a slice of fields
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		question := record[0]
		answer := record[len(record)-1]
		fmt.Println(question, answer)
	}
}
