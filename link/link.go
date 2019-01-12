package link

import (
	"io"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Link represents a link (<a href="..."></a>) in an HTML document.
type Link struct {
	Href string
	Text string
}

// Parse obtains the links present in the given HTML document by parsing it.
// If the given reader is not a valid HTML document, an error is returned,
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse HTML")
	}

	return findLinks(doc), nil
}

// findLinks iterates through the nodes in the HTML document in a depth-first
// (DFS) algorithm, recursively obtaining all the links present in the document.
func findLinks(node *html.Node) []Link {
	if node.Type == html.ElementNode && node.DataAtom == atom.A {
		// We found a link!!
		return []Link{buildLink(node)}
	}

	var links []Link

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, findLinks(c)...)
	}

	return links
}

// buildLink takes an "<a href...>" node and builds a Link structure from it.
// Extracting the "href" attribute and the inner text,
func buildLink(node *html.Node) Link {
	href, _ := findHref(node.Attr)
	text := extractText(node)

	return Link{
		Href: href,
		Text: text,
	}
}

// findHref obtains the value of the "href" attribute from the slice of
// attributes of an HTML node. If the node doesn't have an "href" attribute,
// the default empty string is returned. To check whether this empty string
// is due to a missing argument or to an empty "href", you can check the second
// return argument for presence.
func findHref(attrs []html.Attribute) (string, bool) {
	for _, attr := range attrs {
		if attr.Key == "href" {
			return attr.Val, true
		}
	}

	return "", false
}

// extractText takes an HTML node and returns all the text inside of it properly stripped,
// iterating through the inner elements in a DFS fashion.
func extractText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	} else if node.Type != html.ElementNode {
		return ""
	}

	var text string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}

	return strings.Join(strings.Fields(text), " ")
}
