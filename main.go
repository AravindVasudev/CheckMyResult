package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const aucoe string = "http://aucoe.annauniv.edu/cgi-bin/result/cgrade.pl?regno="

var wg sync.WaitGroup

type student struct {
	RegisterNumber string `json:"registerNumber"`
	EmailID        string `json:"emailID"`
}

type result struct {
	student
	Name       string
	Department string
	Results    map[string]string
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

func requestAUCOE(stud student, results chan result) {
	defer wg.Done()

	doc, err := goquery.NewDocument(aucoe + stud.RegisterNumber)
	if err != nil {
		panic(err)
	}

	var extracted []string
	doc.Find("td[bgcolor='#fffaea']").Each(func(i int, s *goquery.Selection) {
		extracted = append(extracted, s.Text())
	})

	res := result{
		student:    stud,
		Name:       extracted[1],
		Department: extracted[2],
		Results:    make(map[string]string),
	}

	for i := 6; i < len(extracted); i += 3 {
		res.Results[extracted[i]] = extracted[i+1]
	}

	results <- res
}

func main() {
	var students []student
	jsonFromFile("./students.json", &students)

	results := make(chan result, 256)
	for _, stud := range students {
		wg.Add(1)
		go requestAUCOE(stud, results)
	}

	wg.Wait()
	close(results)

	for res := range results {
		fmt.Printf("%v - %v\n", res.student.EmailID, res.Name)
	}
}
