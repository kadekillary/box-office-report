package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/gocolly/colly"
)

const url string = "http://www.boxofficereport.com/trailerviews/trailerviews.html"

var headers = []string{"Rank", "Film", "WeeklyViews", "TotalViews", "ReleaseDate", "TrailerCount"}

type youtubeData struct {
	data [][]string
}

// ToCSV write results from youtubeData to CSV file
func (y *youtubeData) ToCSV(w io.Writer, headers []string) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	writer.Write(headers)
	err := writer.WriteAll(y.data)
	if err != nil {
		log.Fatalf("Failed to write data: %s", err)
	}
	return nil
}

func main() {
	// Used to extract all weekly result links on homepage
	c := colly.NewCollector()

	// Separate collector to grab weekly results
	pageCollector := c.Clone()

	fileName := ExtractFileName(url)
	file, _ := os.Create(fileName)
	defer file.Close()

	c.OnHTML("a.hover2", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		pageCollector.Visit(link)
	})

	youtubeViews := youtubeData{}

	pageCollector.OnHTML("tr", func(e *colly.HTMLElement) {
		temp := make([]string, 6)
		temp[0] = e.ChildText("td.classname1:nth-of-type(1)")
		temp[1] = e.ChildText("td.classname1:nth-of-type(2)")
		temp[2] = e.ChildText("td.classname1:nth-of-type(3)")
		temp[3] = e.ChildText("td.classname1:nth-of-type(4)")
		temp[4] = e.ChildText("td.classname1:nth-of-type(5)")
		temp[5] = e.ChildText("td.classname1:nth-of-type(6)")

		// Ignore empty rows (i.e. any row without a weekly rank)
		if temp[0] != "" {
			youtubeViews.data = append(youtubeViews.data, temp)
		}
	})

	c.Visit(url)

	youtubeViews.ToCSV(file, headers)

}
