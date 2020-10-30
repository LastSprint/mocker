package main

import (
	"net/url"
	"testing"
)

func TestGetStringPathFromHostAndScheme(t *testing.T) {
	// Arrange

	url := url.URL{
		Scheme: "http",
		Host:   "test.host.ex",
	}

	// Act

	res := getDirPathFromURL(&url, "")

	// Assert

	if res != "test.host.ex" {
		t.Fail()
	}
}

func TestGetStringPathFromUrlWithPath(t *testing.T) {
	// Arrange

	url := url.URL{
		Scheme: "http",
		Host:   "test.host.ex",
		Path:   "path/to/endpoint",
	}

	// Act

	res := getDirPathFromURL(&url, "")

	// Assert

	if res != "test.host.ex/path/to/endpoint" {
		t.Fail()
	}
}

func TestGetStringPathFromUrlWithParams(t *testing.T) {
	// Arrange

	url := url.URL{
		Scheme:   "http",
		Host:     "test.host.ex",
		Path:     "path/to/endpoint",
		RawQuery: "isExample=true",
	}

	// Act

	res := getDirPathFromURL(&url, "")

	// Assert

	if res != "test.host.ex/path/to/endpoint" {
		t.Fail()
	}
}

func TestGetUniqNameGenerateReallyUniqNames(t *testing.T) {
	// Arrange

	iterationCount := 10000
	generated := make([]string, iterationCount)

	// Act

	for i := 0; i < iterationCount; i++ {
		generated[i] = getUniqName()
	}

	// Assert

	for i, item := range generated {
		for j, nextItem := range generated {
			if i != j && nextItem == item {
				t.Fail()
			}
		}
	}
}
