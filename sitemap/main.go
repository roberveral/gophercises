package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/roberveral/gophercises/link"
)

const xmlns string = "http://www.sitemaps.org/schemas/sitemap/0.9"

type urlset struct {
	Namespace string `xml:"xmlns,attr"`
	Urls      []loc  `xml:"url"`
}

type loc struct {
	Loc string `xml:"loc"`
}

func main() {
	domain := flag.String("domain", "https://gophercises.com", "Domain to obtain the sitemap for")
	flag.Parse()

	log.Printf("Building SiteMap for domain: %s\n", *domain)

	domainURL, err := url.Parse(*domain)
	if err != nil {
		log.Fatal(err)
		return
	}

	nodes, err := bfs(*domainURL, getAndFilter(*domainURL))
	if err != nil {
		log.Fatal(err)
		return
	}

	sitemap := urlset{Namespace: xmlns}
	for _, node := range nodes {
		sitemap.Urls = append(sitemap.Urls, loc{Loc: node.String()})
	}

	writer := os.Stdout

	fmt.Fprint(writer, xml.Header)
	enc := xml.NewEncoder(writer)
	enc.Indent("", "  ")
	if err := enc.Encode(sitemap); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintln(writer)
}

func bfs(start url.URL, expandFn func(url.URL) ([]url.URL, error)) ([]url.URL, error) {
	queue := []url.URL{start}
	visited := []url.URL{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		children, err := expandFn(node)
		if err != nil {
			return nil, err
		}

		for _, child := range children {
			if !contains(visited, child) && !contains(queue, child) {
				queue = append(queue, child)
			}
		}

		visited = append(visited, node)
	}

	return visited, nil
}

func contains(slice []url.URL, node url.URL) bool {
	for _, elem := range slice {
		if elem == node {
			return true
		}
	}

	return false
}

func getAndFilter(domain url.URL) func(url.URL) ([]url.URL, error) {
	return func(node url.URL) ([]url.URL, error) {
		response, err := http.Get(node.String())
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		links, err := link.Parse(response.Body)
		if err != nil {
			return nil, err
		}

		var filteredLinks []url.URL

		for _, lnk := range links {
			if parsedURL, ok := filterLink(lnk, domain); ok {
				filteredLinks = append(filteredLinks, *parsedURL)
			}
		}

		return filteredLinks, nil
	}
}

func filterLink(lnk link.Link, domain url.URL) (*url.URL, bool) {
	parsedURL, err := url.Parse(lnk.Href)

	if err == nil && parsedURL.Host == "" {
		parsedURL, err = url.Parse(domain.String() + lnk.Href)
	}

	if err == nil && strings.HasPrefix(parsedURL.Scheme, "http") && parsedURL.Host == domain.Host {
		return parsedURL, true
	}

	return nil, false
}
