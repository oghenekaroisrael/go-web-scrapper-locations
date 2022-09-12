package scrapper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type state struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func StateScrapper() {
	count := 0
	states := make([]state, 0)
	// allLocations := make([]locationItem, 0)
	c := colly.NewCollector(
		colly.AllowedDomains("postcode.com.ng", "www.postcode.com.ng"),
		colly.MaxDepth(5),
		// colly.AllowURLRevisit(),
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
		//Delay: 5*time.Second,
	})

	//request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		count += 1
		fmt.Println(count)
	})

	// State
	c.OnHTML("div#pg-22-3 ul li a", func(h *colly.HTMLElement) {
		var astate = h.Text
		link := h.Attr("href")
		h.Request.Ctx.Put("state", astate)
		s := state{
			Name: astate,
			Link: link,
		}
		states = append(states, s)
		// c.Wait()
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished Scraping :", count)
	})

	c.Visit("https://www.postcode.com.ng")
	c.Wait()

	content, err := json.Marshal(states)

	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("states.json", content, 0644)
}
