package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type locationItem2 struct {
	Location string `json:"location"`
	Lga      string `json:"lga"`
	PostCode string `json:"post_code"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

func main() {
	allLocations := make([]locationItem2, 0)
	c := colly.NewCollector(
		colly.AllowedDomains("nigeriazipcodes.com"),
	)
	// State
	c.OnHTML("div#main div#content div.post-content", func(h *colly.HTMLElement) {
		h.ForEach("p strong a", func(i int, h *colly.HTMLElement) {
			var state = h.Text
			if strings.Contains(state, "List Of Towns and Villages in") {
				var link = h.Request.AbsoluteURL(h.Attr("href"))
				absLink := h.Request.AbsoluteURL(link)
				h.Request.Ctx.Put("state", state)
				h.Request.Visit(absLink)
			}

		})
	})

	// LGA in State
	c.OnHTML("div#main div#content div.post-content", func(h *colly.HTMLElement) {
		state := h.Request.Ctx.Get("state")
		h.ForEach("p:first-child strong", func(i int, h *colly.HTMLElement) {
			var lga = h.ChildText("a")
			if strings.Contains(state, "List Of Towns and Villages in") {
				var link = h.Request.AbsoluteURL(h.ChildAttr("a", "href"))
				absLink := h.Request.AbsoluteURL(link)
				h.Request.Ctx.Put("state", state)
				h.Request.Ctx.Put("lga", lga)
				h.Request.Visit(absLink)
			}
		})
	})

	// Towns in STATE LGA
	c.OnHTML("div#main div#content div.post-content", func(h *colly.HTMLElement) {
		state := h.Request.Ctx.Get("state")
		lga := h.Request.Ctx.Get("lga")

		// if page has list
		fmt.Println("have a list")
		h.ForEach("ul", func(i int, h *colly.HTMLElement) {
			h.ForEach("li", func(i int, h *colly.HTMLElement) {
				var alocation = h.Text
				locs := ""
				pc := ""
				var splitLocation []string
				if strings.Contains(alocation, "-") {
					splitLocation = strings.Split(alocation, "-")
				}
				// if strings.Contains(alocation, "=>") {
				// 	splitLocation = strings.Split(alocation, "=>")
				// }
				if !strings.Contains(alocation, "L.G A Zip Codes") && !strings.Contains(alocation, "Town Area Zip Codes") {
					if len(splitLocation) == 2 {
						locs = splitLocation[0]
						pc = splitLocation[1]
					} else {
						locs = alocation
					}

					state_replacer := strings.NewReplacer("Complete List Of Towns, Villages and Zip Codes Of ", "")
					state = state_replacer.Replace(state)

					// lga_replacer := strings.NewReplacer("List of Towns and Villages in", "", " LGA", "")
					// lga = lga_replacer.Replace(state)

					loc := locationItem2{
						State:    state,
						Lga:      lga,
						PostCode: pc,
						Country:  "Nigeria",
						Location: locs,
					}

					// fmt.Printf("Location: %v, State: %v, LGA: %v, PostCode: %v \n", loc.Location, loc.State, loc.Lga, loc.PostCode)
					allLocations = append(allLocations, loc)
				}
			})
		})

		// fmt.Println("I have a table")
		// h.ForEach("table tbody tr", func(i int, h *colly.HTMLElement) {
		// 	var locs = h.ChildText("td:first-child")
		// 	var pc = h.ChildText("td:last-child")
		// 	state_replacer := strings.NewReplacer("Complete List Of Towns, Villages and Zip Codes Of ", "")
		// 	state = state_replacer.Replace(state)

		// 	lga_replacer := strings.NewReplacer("List of Towns and Villages in", "", " LGA", "")
		// 	lga = lga_replacer.Replace(state)

		// 	loc := locationItem2{
		// 		State:    state,
		// 		Lga:      lga,
		// 		PostCode: pc,
		// 		Country:  "Nigeria",
		// 		Location: locs,
		// 	}
		// 	allLocations = append(allLocations, loc)
		// })
	})

	c.Visit("https://nigeriazipcodes.com/5449/list-of-towns-and-villages-in-nigeria-by-states/")

	content, err := json.Marshal(allLocations)

	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("trash.json", content, 0644)
}
