package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type student struct {
	RegisterNumber string `json:"registerNumber"`
	EmailID        string `json:"emailID"`
}

func jsonFromFile(fileName string, store interface{}) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal(raw, &store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func main() {
	var students []student
	jsonFromFile("./students.json", &students)

	for i, stud := range students {
		fmt.Printf("%v : %v\n", i, stud)
	}
}
