package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type local struct {
	Name     string `json:"name"`
	PostCode string `json:"post_code"`
}

type state struct {
	Name string `json:"name"`
	Link string `json:"post_code"`
}
type locationItem struct {
	Location  []local  `json:"place"`
	Lga       []string `json:"lga"`
	States    []state  `json:"state"`
	Countries []string `json:"country"`
}

func main() {
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
	stateCollector := c.Clone()

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
		stateCollector.Visit(h.Request.AbsoluteURL(link))
		// c.Wait()
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("OnScraped:", r.StatusCode)
		fmt.Println("Starting LGA Crawls")

	})

	// LGA in State
	stateCollector.OnHTML("div#main article#post-53 div.entry-content ul li a", func(h *colly.HTMLElement) {
		fmt.Println(h.Request.URL.String())
		state := h.Request.Ctx.Get("state")
		var lga = h.Text
		fmt.Printf("LGA Name: %v \n", lga)
		// var link = h.Attr("href")
		h.Request.Ctx.Put("state", state)
		h.Request.Ctx.Put("lga", lga)
		// fmt.Printf("State Link: %v, LGA link: %v \n\n", state, link)
		// stateCollector.Visit(link)
		// stateCollector.Wait()
	})

	// Towns in STATE LGA
	// c.OnHTML("div#main div#content div.post-content", func(h *colly.HTMLElement) {
	// 	state := h.Request.Ctx.Get("state")
	// 	lga := h.Request.Ctx.Get("lga")
	// 	fmt.Printf("State: %v, LGA: %v \n", state, lga)

	// 	h.ForEach("ul li", func(i int, h *colly.HTMLElement) {
	// 		var alocation = h.Text
	// 		locs := ""
	// 		pc := ""
	// 		var splitLocation []string
	// 		if strings.Contains(alocation, "-") {
	// 			splitLocation = strings.Split(alocation, "-")
	// 		}
	// 		if strings.Contains(alocation, "=>") {
	// 			splitLocation = strings.Split(alocation, "=>")
	// 		}
	// 		if !strings.Contains(alocation, "L.G A Zip Codes") && !strings.Contains(alocation, "Town Area Zip Codes") {
	// 			if len(splitLocation) == 2 {
	// 				locs = splitLocation[0]
	// 				pc = splitLocation[1]
	// 			} else {
	// 				locs = alocation
	// 			}

	// 			state_replacer := strings.NewReplacer("Complete List Of Towns, Villages and Zip Codes Of ", "")
	// 			state = state_replacer.Replace(state)

	// 			lga_replacer := strings.NewReplacer("List of Towns and Villages in", "", " LGA", "")
	// 			lga = lga_replacer.Replace(state)

	// 			loc := locationItem2{
	// 				State:    state,
	// 				Lga:      lga,
	// 				PostCode: pc,
	// 				Country:  "Nigeria",
	// 				Location: locs,
	// 			}

	// 			fmt.Printf("Location: %v, State: %v, LGA: %v, PostCode: %v \n", loc.Location, loc.State, loc.Lga, loc.PostCode)
	// 			allLocations = append(allLocations, loc)
	// 		}
	// 	})
	// })

	// // Before making a request print "Visiting ..."
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL.String())
	// })

	c.Visit("https://www.postcode.com.ng")
	c.Wait()

	// content, err := json.Marshal(allLocations)

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// os.WriteFile("trash.json", content, 0644)
}
