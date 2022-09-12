package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type locationItem struct {
	Location string `json:"location"`
	Lga      string `json:"lga"`
	PostCode string `json:"post_code"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

func mainw() {
	allLocations := make([]locationItem, 0)
	c := colly.NewCollector(
		colly.AllowedDomains("nigeriazipcodes.com"),
	)

	c.OnHTML("div#main div#content div.post-content", func(h *colly.HTMLElement) {
		h.ForEach("p strong a", func(i int, h *colly.HTMLElement) {
			var state = h.Text
			var link = h.Request.AbsoluteURL(h.Attr("href"))
			// fmt.Println(state)
			// fmt.Println(link)
			absLink := h.Request.AbsoluteURL(link)
			h.Request.Visit(absLink)
			c.OnHTML("div#main div#content div.post-content", func(h *colly.HTMLElement) {
				h.ForEach("ul li a", func(i int, h *colly.HTMLElement) {
					var lga = h.Text
					var link = h.Request.AbsoluteURL(h.Attr("href"))
					// fmt.Println(state)
					// fmt.Println(link)
					absLink := h.Request.AbsoluteURL(link)
					h.Request.Visit(absLink)

					c.OnHTML("div#main div#content div.post-content", func(h *colly.HTMLElement) {
						h.ForEach("table tbody tr", func(i int, h *colly.HTMLElement) {
							var locs = h.ChildText("td:first-child")
							var pc = h.ChildText("td:last-child")
							loc := locationItem{
								State:    state,
								Lga:      lga,
								PostCode: pc,
								Country:  "Nigeria",
								Location: locs,
							}
							allLocations = append(allLocations, loc)
						})
					})
				})
			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		fmt.Println(r)
	})

	c.Visit("https://nigeriazipcodes.com/5449/list-of-towns-and-villages-in-nigeria-by-states/")

	content, err := json.Marshal(allLocations)

	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("trash.json", content, 0644)
}
