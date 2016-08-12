package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// define our identifiers for letter- and item-urls
var pageUrlIdentifier string = "li.glossar_link.buchstaben_navi_has_entry"
var itemUrlIdentifier string = "li.glossar_link:not(.buchstaben_navi_has_entry)"

// FetchGlossaryItemsFrom will fetch all glossary items found on the given kofam site
func FetchGlossaryItemsFrom(kofamUrl string) ([]GlossaryItem, error) {

	// parse the url to find the base
	urlScheme, _ := url.Parse(kofamUrl)
	baseUrl := urlScheme.Scheme + "://" + urlScheme.Host

	glossaryLetterPages := []GlossaryPage{}
	glossaryItems := []GlossaryItem{}

	// go to the swissethics page and load all documents
	doc, err := goquery.NewDocument(kofamUrl)
	if err != nil {
		return nil, err
	}

	fmt.Println("Checking:", kofamUrl)

	// extract the urls to load each letter
	doc.Find(pageUrlIdentifier).Each(func(i int, s *goquery.Selection) {

		// check if the link has an onClick attribute
		link, exists := s.Attr("onclick")

		if exists == true {
			// define a new glossary page
			page := GlossaryPage{}

			// extract the letter from the text property
			page.Letter = s.Text()

			// we are only looking for the url
			page.Url = prefixWithBaseUrl(baseUrl, extractLinkOnly(link))

			// add the links to all our pages
			glossaryLetterPages = append(glossaryLetterPages, page)
		}
	})

	// make sure that all links are absolute
	for index := range glossaryLetterPages {

		// fetch all links on the glossar page
		itemsOnPage, err := fetchItemsOnPage(glossaryLetterPages[index], baseUrl)
		if err != nil {
			return nil, err
		}

		// append the items on this page
		glossaryItems = append(glossaryItems, itemsOnPage...)
	}

	// load the detail information for each item
	for index, item := range glossaryItems {
		glossaryItems[index] = fetchDetails(item)
	}

	// return all definitions
	return glossaryItems, nil
}

// fetchItemsOnPage will return the urls of all items on the given page
func fetchItemsOnPage(page GlossaryPage, baseUrl string) ([]GlossaryItem, error) {

	items := []GlossaryItem{}

	// go to the swissethics page and load all documents
	doc, err := goquery.NewDocument(page.Url)
	if err != nil {
		return nil, err
	}

	// extract the urls to load each letter
	doc.Find(itemUrlIdentifier).Each(func(i int, s *goquery.Selection) {

		// check if the link has an onClick attribute
		link, exists := s.Attr("onclick")

		if exists == true {
			// define a new glossary item
			item := GlossaryItem{}

			// extract the letter from the text property
			item.Name = s.Text()

			// we are only looking for the url
			item.Url = prefixWithBaseUrl(baseUrl, extractLinkOnly(link))

			// add the links to all our items
			items = append(items, item)
		}
	})

	return items, nil
}

// fetchDetails will fetch the details for a given glossary item
func fetchDetails(item GlossaryItem) GlossaryItem {

	// return the item if url was not specified
	if item.Url == "" {
		return item
	}

	// go to the swissethics page and load all documents
	doc, err := goquery.NewDocument(item.Url)
	if err != nil {
		return item
	}

	// title := doc.Find("h4").First().Text()
	subtitle := doc.Find("h5").First().Text()
	description := doc.Find("p.no_margin").First().Text()

	item.English = doc.Find("h4").First().Next().First().Text()
	item.Source = doc.Find("p.quellenangabe").First().Text()

	// extract the full content of the element
	doc.Find(":nth-child(n+3)").Each(func(i int, s *goquery.Selection) {
		item.ContentFull = item.ContentFull + "\n" + s.Text()
	})
	item.ContentFull = strings.Trim(item.ContentFull, " \n")

	if subtitle != "" {
		item.Description = subtitle + " " + description
	} else {
		item.Description = description
	}

	return item
}

// --- HELPER FUNCTIONS ---

// extractLinkOnly will extract the url from the onClick javascript handler: foo('[url]')
func extractLinkOnly(url string) string {
	url = strings.Replace(url, "foo('", "", -1)
	url = strings.Replace(url, "')", "", -1)
	return url
}

// prefixWithBaseUrl will make sure that all links are absolute
func prefixWithBaseUrl(base string, url string) string {
	// return the url if it contains the base url
	if strings.Contains(url, base) == true {
		return url
	}

	// return the url if it starts with http
	if strings.HasPrefix(url, "http") {
		return url
	}

	if strings.HasSuffix(base, "/") == false && strings.HasPrefix(url, "/") == false {
		base = base + "/"
	}

	return base + url
}
