package link

import (
	"os"
	"testing"
)

func fileParseTest(t *testing.T, testFile string, expected []Link) {
	file, err := os.Open(testFile)
	if err != nil {
		t.Error("Missing test file: ", testFile)
		return
	}

	links, err := Parse(file)
	if err != nil {
		t.Errorf("Expected valid result, but an error was returned: %+v", err)
		return
	}

	if len(links) != len(expected) {
		t.Errorf("Expected a result with %d elements, but got %d", len(expected), len(links))
		return
	}

	for i, link := range links {
		if link != expected[i] {
			t.Errorf("Expected result list: %+v, but got: %+v", expected, links)
			return
		}
	}
}

func TestParseEx1(t *testing.T) {
	testFile := "ex1.html"
	expected := []Link{
		Link{"/other-page", "A link to another page"},
	}

	fileParseTest(t, testFile, expected)
}

func TestParseEx2(t *testing.T) {
	testFile := "ex2.html"
	expected := []Link{
		Link{"https://www.twitter.com/joncalhoun", "Check me out on twitter"},
		Link{"https://github.com/gophercises", "Gophercises is on Github!"},
	}

	fileParseTest(t, testFile, expected)
}

func TestParseEx3(t *testing.T) {
	testFile := "ex3.html"
	expected := []Link{
		Link{"#", "Login"},
		Link{"/lost", "Lost? Need help?"},
		Link{"https://twitter.com/marcusolsson", "@marcusolsson"},
	}

	fileParseTest(t, testFile, expected)
}

func TestParseEx4(t *testing.T) {
	testFile := "ex4.html"
	expected := []Link{
		Link{"/dog-cat", "dog cat"},
	}

	fileParseTest(t, testFile, expected)
}
