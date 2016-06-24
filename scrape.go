package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func readLines(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")
	return lines[:len(lines)-1]
}

func findArtistPage(url string, artist string) (artist_page string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	// Break out the loop by returning false once we've found the corresponding link
	doc.Find(".artists_index_list li a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.EqualFold(s.Text(), artist) {
			artist_page, _ = s.Attr("href")
			return false
		}
		return true
	})
	return
}

// [TODO] Instead of returnign a single list of urls I need to return a data
// type that returns the album name, along with corresponding url i.e. {'album_name': url}
func findAlbumPages(url string) (album_pages []string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".album_list .album_link").Each(func(i int, s *goquery.Selection) {
		album_page, _ := s.Attr("href")
		album_pages = append(album_pages, album_page)
	})
	return
}

// [TODO] Do same here
func findSongPages(url string) (song_pages []string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".album_tracklist ul li a").Each(func(i int, s *goquery.Selection) {
		song_page, _ := s.Attr("href")
		song_pages = append(song_pages, song_page)
	})
	return
}

func main() {
	index_url := "http://genius.com/artists-index/"
	artists := readLines("artists.txt")
	for _, artist := range artists {
		first_char := artist[0:1]
		page_url := findArtistPage(index_url+first_char, artist)
		if len(page_url) != 0 {
			fmt.Printf("\n%s\n", page_url)
			album_pages := findAlbumPages(page_url)
			for _, album_page := range album_pages {
				fmt.Println(album_page)
				song_pages := findSongPages("http://genius.com" + album_page)
				for _, song_page := range song_pages {
					fmt.Println(song_page)
				}
			}
		} else {
			fmt.Printf("Artist not found: %s\n", artist)
		}
	}
}
