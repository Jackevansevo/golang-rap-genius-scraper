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

// For finding songs
// findPages(root_url+album_page.url, ".album_tracklist ul li a .song_title")

type artistPage struct {
	artist_name, url string
}

type albumPage struct {
	artist_name, album_name, url string
}

type songPage struct {
	artist_name, album_name, song_name, url string
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

func findArtistPage(artist string, ch chan<- artistPage) {
	index_url := "http://genius.com/artists-index/"
	doc := scrapePage(index_url + artist[0:1])
	var artist_page string = ""
	doc.Find(".artists_index_list li a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.EqualFold(s.Text(), artist) {
			artist_page, _ = s.Attr("href")
			return false
		}
		return true
	})
	ch <- artistPage{artist, artist_page}
}

func findAlbumPages(url string, artist string, html string, ch chan<- []albumPage) {
	doc := scrapePage(url)
	pages := []albumPage{}
	doc.Find(html).Each(func(i int, s *goquery.Selection) {
		page_url, _ := s.Attr("href")
		page_name := strings.TrimSpace(s.Text())
		pages = append(pages, albumPage{artist, page_name, page_url})
	})
	ch <- pages
}

func main() {
	artists := readLines("artists.txt")
	artistPageCh := make(chan artistPage)
	for _, artist := range artists {
		go findArtistPage(artist, artistPageCh)
	}

	artist_pages := []artistPage{}
	for i := 0; i < len(artists); i++ {
		artist_pages = append(artist_pages, <-artistPageCh)
	}

	albumPageCh := make(chan []albumPage)
	for elem := range artist_pages {
		url := artist_pages[elem].url
		name := artist_pages[elem].artist_name
		go findAlbumPages(url, name, ".album_list .album_link", albumPageCh)
	}

	album_pages := [][]albumPage{}
	for i := 0; i < len(artist_pages); i++ {
		album_pages = append(album_pages, <-albumPageCh)
	}

	// [TODO] Figure out a way of linking together the structs
	// artists -> album_pages -> songs
	for artist := range album_pages {
		for album := range album_pages[artist] {
			fmt.Println(album_pages[artist][album])
		}
	}
}
