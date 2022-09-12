package main

import "github.com/oghenekaroisrael/go-web-scrapper-locations/scrapper"

func main() {
	// scrapper.StateScrapper()
	// scrapper.Crawl_The_State("http://postcode.com.ng/abia-state-lga-nigeria-postcode/", "Abia")
	scrapper.Crawl_The_Country("http://postcode.com.ng/", "Nigeria")
}
