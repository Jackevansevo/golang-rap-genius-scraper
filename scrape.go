package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// [TODO] Fix script to look up artists who's name begins with a number, e.g.
// ["2 Chainz", "2Pac", "50 Cent"]

type urlName struct {
	name, url string
}

func readLines(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")
	return lines[:len(lines)-1]
}

func scrapePage(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func findArtistPage(url string, artist string) (artist_page string) {
	doc := scrapePage(url)
	doc.Find(".artists_index_list li a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.EqualFold(s.Text(), artist) {
			artist_page, _ = s.Attr("href")
			return false
		}
		return true
	})
	return
}

func findPages(url string, html string) *[]urlName {
	doc := scrapePage(url)
	pages := []urlName{}
	doc.Find(html).Each(func(i int, s *goquery.Selection) {
		page_url, _ := s.Attr("href")
		page_name := strings.TrimSpace(s.Text())
		pages = append(pages, urlName{page_name, page_url})
	})
	return &pages
}

func main() {
	index_url := "http://genius.com/artists-index/"
	artists := readLines("artists.txt")
	// [TODO] Don't loop through artists, map through instead
	for _, artist := range artists {
		first_char := artist[0:1]
		artist_page_url := findArtistPage(index_url+first_char, artist)
		if len(artist_page_url) != 0 {
			fmt.Printf("\n%s\n", artist_page_url)
			album_pages := findPages(artist_page_url, ".album_list .album_link")
			for _, album_page := range *album_pages {
				fmt.Printf("\n%s\n", album_page.name)
				full_album_page_url := "http://genius.com" + album_page.url
				song_pages := findPages(full_album_page_url, ".album_tracklist ul li a")
				for _, song_page := range *song_pages {
					fmt.Println(song_page.name)
				}
			}
		} else {
			fmt.Printf("Artist not found: %s\n", artist)
		}
	}
}
