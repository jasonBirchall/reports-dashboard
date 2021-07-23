package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

var (
	userGuide   = flag.String("userGuide", "user-guide.cloud-platform.service.justice.gov.uk", "Full URL of the userguide.")
	runBook     = flag.String("runBook", "runbooks.cloud-platform.service.justice.gov.uk", "Full URL of the runbook site.")
	currentTime = time.Now()
)

func main() {
	flag.Parse()

	m := collect()

	jsonString, err := json.Marshal(m)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(jsonString)
}

func collect() map[string]string {
	// spider url looking for links to other pages
	// return a hash of pages: { pageUrl : needsReview? }
	c := colly.NewCollector(
		colly.AllowedDomains(*userGuide, *runBook),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "justice",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// Look for div value "data-last-reviewed-on" which contains an int value
	expired := make(map[string]string)
	c.OnHTML("div[data-last-reviewed-on]", func(e *colly.HTMLElement) {
		lastReviewed, _ := e.DOM.Attr("data-last-reviewed-on")
		page := e.Request.URL.String()
		if lastReviewed < currentTime.Format("2006-01-02") {
			expired[page] = lastReviewed
		}
	})

	c.Visit("https://" + *userGuide)
	c.Visit("https://" + *runBook)

	c.Wait()

	return expired
}
