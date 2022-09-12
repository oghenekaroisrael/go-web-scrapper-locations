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
type local struct {
	Name     string `json:"name"`
	PostCode string `json:"post_code"`
}

type State struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type StateItem struct {
	CountryName string  `json:"countryName"`
	StateName   string  `json:"stateName"`
	LGA         []lga   `json:"lga"`
	Places      []local `json:"places"`
}

type lga struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type States struct {
	States []State `json:"states"`
}

type locationItem struct {
	Location  []local  `json:"place"`
	Lga       []lga    `json:"lga"`
	States    []state  `json:"state"`
	Countries []string `json:"country"`
}

func LgaScrapper() {
	// count := 0

	states := make([]state, 0)

	lgas := make([]lga, 0)

	streets := make([]local, 0)

	countries := make([]string, 0)

	// // Open our jsonFile
	// jsonFile, err := os.Open("../states.json")
	// // if we os.Open returns an error then handle it
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Successfully Opened States.json")
	// // defer the closing of our jsonFile so that we can parse it later on
	// defer jsonFile.Close()

	// // read our opened jsonFile as a byte array.
	// byteValue, _ := ioutil.ReadAll(jsonFile)

	// // we initialize our Users array
	// var states States

	// // we unmarshal our byteArray which contains our
	// // jsonFile's content into 'users' which we defined above
	// json.Unmarshal(byteValue, &states)

	// // max := len(states.States) - 1

	// // allLocations := make([]locationItem, 0)
	// c := colly.NewCollector(
	// 	colly.AllowedDomains("postcode.com.ng", "www.postcode.com.ng"),
	// 	colly.MaxDepth(5),
	// 	// colly.AllowURLRevisit(),
	// 	colly.Async(true),
	// 	// colly.Debugger(&debug.LogDebugger{}),
	// )
	// c.Limit(&colly.LimitRule{
	// 	DomainGlob:  "*",
	// 	Parallelism: 4,
	// 	//Delay: 5*time.Second,
	// })

	// //request
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL.String())
	// 	count += 1
	// 	fmt.Println(count)
	// })

	// // we iterate through every user within our users array and
	// // print out the user Type, their name, and their facebook url
	// // as just an example
	// // for i := 0; i < len(states.States); i++ {
	// // 	// LGA in State
	// c.OnHTML("div#main div.corp-content-wrapper div.entry-content ul li a", func(h *colly.HTMLElement) {
	// 	var lg = h.Text
	// 	var link = h.Attr("href")
	// 	fmt.Printf("LGA name: %v, LGA link: %v \n\n", lg, link)
	// 	l := lga{
	// 		Name: lg,
	// 		Link: link,
	// 	}
	// 	lgas = append(lgas, l)

	// })
	// // }

	// c.OnScraped(func(r *colly.Response) {
	// 	fmt.Println("Finished Scraping :", count)
	// })

	// c.Visit("http://postcode.com.ng/abia-state-lga-nigeria-postcode/")
	// c.Wait()

	// content, err := json.Marshal(lgas)

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// os.WriteFile("lgas.json", content, 0644)
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("postcode.com.ng", "www.postcode.com.ng"),
		colly.Async(true),
	)

	// State
	c.OnHTML("div#pg-22-3 div.siteorigin-widget-tinymce ul li", func(h *colly.HTMLElement) {
		var astate = h.ChildText("a")
		link := h.ChildAttr("a", "href")
		s := state{
			Name: astate,
			Link: link,
		}
		countries = append(countries, "Nigeria")
		states = append(states, s)
		fmt.Println(s)
		c.Visit(h.Request.AbsoluteURL(link))
	})

	// LGA
	c.OnHTML("div#content div.corp-content-wrapper div.entry-content ul li", func(h *colly.HTMLElement) {
		var alga = h.ChildText("a")
		link := h.ChildAttr("a", "href")
		s := lga{
			Name: alga,
			Link: link,
		}
		lgas = append(lgas, s)
		fmt.Println(s)
		c.Visit(h.Request.AbsoluteURL(link))
	})

	// Street and PostCode
	c.OnHTML("div#content div.corp-content-wrapper div.entry-content table.wp-block-table tbody tr", func(h *colly.HTMLElement) {
		var street = h.ChildText("td:first-child")
		var postcode = h.ChildText("td:last-child")

		s := local{
			Name:     street,
			PostCode: postcode,
		}
		streets = append(streets, s)
		fmt.Println("Location ", s)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// c.OnScraped(func(r *colly.Response) {
	// 	fmt.Println("Finished Scraping :", lgas)
	// })

	c.Visit("https://www.postcode.com.ng/")
	c.Wait()

	l := locationItem{
		Lga:       lgas,
		States:    states,
		Countries: countries,
		Location:  streets,
	}

	content, err := json.Marshal(l)

	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("abia.json", content, 0644)
}
