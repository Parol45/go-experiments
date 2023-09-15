package main

import (
	"log"
	"magnet-parser/lib"
	"os"
	"strings"
)

func inArray(val string, array []string) (exists bool) {
	exists = false
	for _, v := range array {
		if val == v {
			exists = true
			return
		}
	}
	return
}

func appendLineToFile(magnet string) {
	f, err := os.OpenFile("magnets.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	if _, err := f.Write([]byte(magnet + "\n\n")); err != nil {
		log.Fatal(err)
		return
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	os.Remove("magnets.txt")
	links := []string{"https://igg-games.ru/"}
	index := 0

	for index < len(links) {
		url := links[index]
		urlOrigin := lib.ParseOriginUrl(url)
		html := lib.GetHtmlByUrl(url)

		// add all found links to queue
		newLinks := lib.ParseAllLinks(html)
		for _, newLink := range newLinks {
			if strings.HasPrefix(newLink, "http") {
				if !inArray(newLink, links) {
					links = append(links, newLink)
				}
			} else if strings.HasPrefix(newLink, "/") && urlOrigin != "" {
				if !inArray(newLink, links) {
					links = append(links, urlOrigin+newLink)
				}
			} else if strings.HasPrefix(newLink, "magnet:?xt=urn:") {
				appendLineToFile(newLink)
			}
		}
		index++
	}

	log.Println("Done.")
}
