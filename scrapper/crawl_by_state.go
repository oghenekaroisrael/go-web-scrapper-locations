package scrapper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type placeItem struct {
	Name     string `json:"name"`
	PostCode string `json:"post_code"`
}

type areaItem struct {
	Area   string      `json:"area"`
	Places []placeItem `json:"places"`
}

type StateItem struct {
	StateName string    `json:"stateName"`
	LGA       []lgaItem `json:"lga"`
}

type CountryItem struct {
	CountryName string      `json:"countryName"`
	States      []StateItem `json:"states"`
}

type lgaItem struct {
	LgaName string     `json:"lgaName"`
	Areas   []areaItem `json:"areas"`
}

type State struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

// Crawls A Local Government For Areas then craw those areas to get places In that Area
func Crawl_The_LGA(url string) []areaItem {
	count := 0
	this_area := areaItem{}
	areas := make([]areaItem, 0)

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only these domains
		colly.AllowedDomains("postcode.com.ng", "www.postcode.com.ng"),
		colly.MaxDepth(1),
		// colly.Async(true),
	)
	c.AllowURLRevisit = false

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Area : ", r.URL.String())
		count += 1
	})

	// Areas
	c.OnHTML("div#primary div.entry-content", func(h *colly.HTMLElement) {

		h.ForEach("h2", func(j int, e *colly.HTMLElement) {
			var name = e.Text
			this_area.Area = name
			fmt.Println(name)

			streets := make([]placeItem, 0)
			h.ForEach("div.entry-content table.wp-block-table tbody tr", func(i int, h *colly.HTMLElement) {
				var street = h.ChildText("td:first-child")
				var postcode = h.ChildText("td:last-child")

				s := placeItem{
					Name:     street,
					PostCode: postcode,
				}
				streets = append(streets, s)
			})
			this_area.Places = streets
			streets = nil

			areas = append(areas, this_area)
		})

	})

	c.Visit(url)
	return areas
}

// Crawls a state for local governments
func Crawl_The_State(url string, name string) []StateItem {
	states := make([]StateItem, 0)

	this_state := StateItem{}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only these domains
		colly.AllowedDomains("postcode.com.ng", "www.postcode.com.ng"),
		// colly.Async(true),
		colly.MaxDepth(1),
	)
	c.AllowURLRevisit = false

	// Crawl A State
	c.OnHTML("div#primary div.entry-content", func(h *colly.HTMLElement) {
		// State Details
		var name = h.ChildText("h1.entry-title")
		fmt.Println(name)
		this_state.StateName = name
		lgas := make([]lgaItem, 0)

		h.ForEach("div.corp-content-wrapper div.entry-content ul li", func(i int, h *colly.HTMLElement) {
			var alga = h.ChildText("a")
			link := h.ChildAttr("a", "href")
			s := lgaItem{}
			s.LgaName = alga
			s.Areas = Crawl_The_LGA(link)
			lgas = append(lgas, s)
		})
		this_state.LGA = lgas
		lgas = nil
		states = append(states, this_state)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting LGA : ", r.URL.String())
	})

	c.Visit(url)
	content, err := json.Marshal(states)

	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("states.json", content, 0644)
	return states
}

// Crawls  a country to get states
func Crawl_By_Country(url string, name string) []CountryItem {

	countries := make([]CountryItem, 0)

	country := CountryItem{}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only these domains
		colly.AllowedDomains("postcode.com.ng", "www.postcode.com.ng"),
		// colly.Async(true),
		colly.MaxDepth(1),
	)
	c.AllowURLRevisit = false

	// Crawl A State
	c.OnHTML("div#primary div.entry-content", func(h *colly.HTMLElement) {
		// State Details
		country.CountryName = "Nigeria"
		h.ForEach("div.corp-content-wrapper div.entry-content ul li", func(i int, h *colly.HTMLElement) {
			var stateName = h.ChildText("a")
			link := h.ChildAttr("a", "href")
			s := CountryItem{}
			s.States = Crawl_The_State(link, stateName)
			// country.
		})
		// countries = append(countries, s)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting LGA : ", r.URL.String())
	})

	c.Visit(url)
	return countries
}
