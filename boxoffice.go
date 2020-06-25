package main

import (
	"encoding/csv"
	"io"
	"os"
	"reflect"

	"github.com/gocolly/colly"
)

const url string = "http://www.boxofficereport.com/trailerviews/trailerviews.html"

var headers = []string{"Rank", "Film", "WeeklyViews", "TotalViews", "ReleaseDate", "TrailerCount"}

type viewsData struct {
	Rank         string
	Film         string
	WeeklyViews  string
	TotalViews   string
	ReleaseDate  string
	TrailerCount string
}

// Use to check whether a structure is empty or not
func (v viewsData) IsStructureEmpty() bool {
	return reflect.DeepEqual(v, viewsData{})
}

type youtubeData []viewsData

// ToCSV writes data from struct to file
func (y *youtubeData) ToCSV(w io.Writer, headers []string) error {
	writer := csv.NewWriter(w)
	writer.Write(headers)
	for _, m := range *y {
		writer.Write([]string{
			m.Rank,
			m.Film,
			m.WeeklyViews,
			m.TotalViews,
			m.ReleaseDate,
			m.TrailerCount,
		})
	}
	writer.Flush()
	return writer.Error()
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
		temp := viewsData{}
		temp.Rank = e.ChildText("td.classname1:nth-of-type(1)")
		temp.Film = e.ChildText("td.classname1:nth-of-type(2)")
		temp.WeeklyViews = e.ChildText("td.classname1:nth-of-type(3)")
		temp.TotalViews = e.ChildText("td.classname1:nth-of-type(4)")
		temp.ReleaseDate = e.ChildText("td.classname1:nth-of-type(5)")
		temp.TrailerCount = e.ChildText("td.classname1:nth-of-type(6)")

		// Don't write any results where the structure is empty
		if temp.IsStructureEmpty() != true {
			youtubeViews = append(youtubeViews, temp)
		}
	})

	c.Visit(url)

	youtubeViews.ToCSV(file, headers)

}
