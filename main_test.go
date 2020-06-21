package main

import "testing"

func TestExtractFileName(t *testing.T) {
	correctFileName := "trailerviews.csv"
	fileName := ExtractFileName(url)
	if fileName != correctFileName {
		t.Errorf("fileName was incorrect, got: %s, wanted: %s", fileName, correctFileName)
	}
}

func TestIsStructureEmpty(t *testing.T) {
	viewsDataTest := viewsData{}
	if viewsDataTest.IsStructureEmpty() != true {
		t.Error("Structure is not empty!")
	}
}
