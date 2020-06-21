package main

import (
	"encoding/csv"
	"os"
	"reflect"
	"strings"

	"github.com/gocolly/colly"
)

const (
	url string = "http://www.boxofficereport.com/trailerviews/trailerviews.html"
)

type viewsData struct {
	Rank         string
	Film         string
	WeeklyViews  string
	TotalViews   string
	ReleaseDate  string
	TrailerCount string
}

// Use to check whether a sturcture is empty or not
func (v viewsData) IsStructureEmpty() bool {
	return reflect.DeepEqual(v, viewsData{})
}

// Grab filename out of URL
// Could refactor to write each week's results into a seperate file
func extractFileName(url string) string {
	splitURL := strings.Split(url, "/")
	extractedFileName := splitURL[len(splitURL)-1]
	fileName := strings.Replace(extractedFileName, "html", "csv", -1)
	return fileName
}

func main() {
	// Used to extract all weekly result links on homepage
	c := colly.NewCollector()

	// Seperate collector to grab weekly results
	pageCollector := c.Clone()

	fileName := extractFileName(url)
	file, _ := os.Create(fileName)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"Rank", "Film", "WeeklyViews", "TotalViews", "ReleaseDate", "TrailerCount"})

	youtubeViews := []viewsData{}

	c.OnHTML("a.hover2", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		pageCollector.Visit(link)
	})

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

	for i := range youtubeViews {
		writer.Write([]string{
			youtubeViews[i].Rank,
			youtubeViews[i].Film,
			youtubeViews[i].WeeklyViews,
			youtubeViews[i].TotalViews,
			youtubeViews[i].ReleaseDate,
			youtubeViews[i].TrailerCount,
		})
	}
}
