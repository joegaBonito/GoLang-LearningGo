package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

type Writer struct {
	Comma   rune // Field delimiter (set to ',' by NewWriter)
	UseCRLF bool // True to use \r\n as the line terminator
	w       *bufio.Writer
}

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50"

var extractedJobChannel = make(chan extractedJob)
var jobsChannel = make(chan []extractedJob)
var writeChannel = make(chan error)

func main() {
	var jobs []extractedJob
	// jobsChannel := make(chan []extractedJob)
	totalPages := getPages()
	fmt.Println(totalPages)
	for i := 0; i < totalPages; i++ {
		go getPage(i, jobsChannel)
	}
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-jobsChannel
		jobs = append(jobs, extractedJobs...)
	}

	fmt.Println(jobs)
	writeJobs(jobs)
	fmt.Println("Finished Extracting Jobs Online")
}

func writeEachRow(w *csv.Writer, jobSlice []string, writeChannel chan<- error) {
	jwErr := w.Write(jobSlice)
	writeChannel <- jwErr
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		go writeEachRow(w, jobSlice, writeChannel)
	}
	for i := 0; i < len(jobs); i++ {
		jwErr := <-writeChannel
		checkErr(jwErr)
	}
}

func getPage(page int, jobsChannel chan<- []extractedJob) {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, extractedJobChannel)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-extractedJobChannel
		jobs = append(jobs, job)
	}

	jobsChannel <- jobs
}

func extractJob(card *goquery.Selection, extractedJobChannel chan<- extractedJob) {
	id, _ := card.Attr("data-jk")

	title := cleanString(card.Find(".title>a").Text())

	location := cleanString(card.Find(".sjcl").Text())

	salary := cleanString(card.Find(".salaryText").Text())

	summary := cleanString(card.Find("").Text())

	extractedJobChannel <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() //closes response to prevent memory leaks
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// fmt.Println(s.Html())
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with status: ", res.StatusCode)
	}
}
