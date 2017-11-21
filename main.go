package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
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

type emailData struct {
	EmailID  string `json:"emailID"`
	Password string `json:"password"`
	Server   string `json:"server"`
}

var emailAuth smtp.Auth
var emailAuthData emailData
var emailTemplate *template.Template

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

	log.Println(res)

	results <- res

	wg.Add(1)
	go sendResultEmail(results)
}

func sendResultEmail(results chan result) {
	defer wg.Done()

	res := <-results

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Semester Result\n"

	buf := new(bytes.Buffer)
	err := emailTemplate.Execute(buf, res)
	if err != nil {
		panic(err)
	}

	msg := []byte(subject + mime + "\n" + buf.String())

	err = smtp.SendMail(emailAuthData.Server+":587", emailAuth, emailAuthData.EmailID, []string{res.student.EmailID}, msg)
	if err != nil {
		panic(err)
	}
}

func main() {
	var students []student
	jsonFromFile("./students.json", &students)

	jsonFromFile("./email_smtp.json", &emailAuthData)

	emailAuth = smtp.PlainAuth(
		"",
		emailAuthData.EmailID,
		emailAuthData.Password,
		emailAuthData.Server,
	)

	var err error
	emailTemplate, err = template.ParseFiles("./email_template.html")
	if err != nil {
		panic(err)
	}

	results := make(chan result, 256)
	for _, stud := range students {
		wg.Add(1)
		go requestAUCOE(stud, results)
	}

	wg.Wait()
	close(results)

	fmt.Println("Done.")
}
