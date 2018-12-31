package urlshort

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYAMLReturnsTheMappings(t *testing.T) {
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	expected := []pathMapping{
		pathMapping{"/urlshort", "https://github.com/gophercises/urlshort"},
		pathMapping{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}

	result, err := parseYAML([]byte(yaml))

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParseYAMLReturnsErrorIfMalformed(t *testing.T) {
	yaml := `{ "path": "/urlshort" }`

	_, err := parseYAML([]byte(yaml))

	assert.NotNil(t, err)
}

func TestParseJSONReturnsTheMappings(t *testing.T) {
	json := `[
		{ "path": "/urlshort", "url": "https://github.com/gophercises/urlshort" },
		{ "path": "/urlshort-final", "url": "https://github.com/gophercises/urlshort/tree/solution" }
	]`

	expected := []pathMapping{
		pathMapping{"/urlshort", "https://github.com/gophercises/urlshort"},
		pathMapping{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}

	result, err := parseJSON([]byte(json))

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestParseJSONReturnsErrorIfMalformed(t *testing.T) {
	json := `path: /urlshort`

	_, err := parseJSON([]byte(json))

	assert.NotNil(t, err)
}

func TestBuildMapReturnsMapIndexedByPath(t *testing.T) {
	pathMappings := []pathMapping{
		pathMapping{"/urlshort", "https://github.com/gophercises/urlshort"},
		pathMapping{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}

	expected := map[string]string{
		"/urlshort":       "https://github.com/gophercises/urlshort",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}

	result := buildMap(pathMappings)

	assert.Equal(t, expected, result)
}

func TestParseAndBuildMapHandlerReturnsHandlerFunc(t *testing.T) {
	inputData := "some test data"
	pathMappings := []pathMapping{
		pathMapping{"/urlshort", "https://github.com/gophercises/urlshort"},
		pathMapping{"/urlshort-final", "https://github.com/gophercises/urlshort/tree/solution"},
	}
	dataParser := func(data []byte) ([]pathMapping, error) {
		assert.Equal(t, []byte(inputData), data)
		return pathMappings, nil
	}

	result, err := parseAndBuildMapHandler(dataParser, []byte(inputData), nil)

	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestParseAndBuildMapHandlerReturnsErrorIfCannotParse(t *testing.T) {
	inputData := "some test data"
	dataParser := func(data []byte) ([]pathMapping, error) {
		assert.Equal(t, []byte(inputData), data)
		return nil, errors.New("error")
	}

	result, err := parseAndBuildMapHandler(dataParser, []byte(inputData), nil)

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestMapHandlerRedirectsToUrlWhenPathFound(t *testing.T) {
	pathsToUrls := map[string]string{
		"/urlshort":       "https://github.com/gophercises/urlshort",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}
	fallback := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	request := httptest.NewRequest("GET", "/urlshort", nil)
	response := httptest.NewRecorder()

	handler := MapHandler(pathsToUrls, fallback)

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "https://github.com/gophercises/urlshort", response.HeaderMap.Get("Location"))
}

func TestMapHandlerForwardsToFallbackWhenPathNotFound(t *testing.T) {
	pathsToUrls := map[string]string{
		"/urlshort":       "https://github.com/gophercises/urlshort",
		"/urlshort-final": "https://github.com/gophercises/urlshort/tree/solution",
	}
	fallback := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	request := httptest.NewRequest("GET", "/missing", nil)
	response := httptest.NewRecorder()

	handler := MapHandler(pathsToUrls, fallback)

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, []byte("Hello World"), body)
}

func TestYAMLHandlerRedirectsToUrlWhenPathFound(t *testing.T) {
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	fallback := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	request := httptest.NewRequest("GET", "/urlshort", nil)
	response := httptest.NewRecorder()

	handler, err := YAMLHandler([]byte(yaml), fallback)

	assert.Nil(t, err)

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "https://github.com/gophercises/urlshort", response.HeaderMap.Get("Location"))
}

func TestJSONHandlerRedirectsToUrlWhenPathFound(t *testing.T) {
	json := `[
		{ "path": "/urlshort", "url": "https://github.com/gophercises/urlshort" },
		{ "path": "/urlshort-final", "url": "https://github.com/gophercises/urlshort/tree/solution" }
	]`
	fallback := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World"))
	})

	request := httptest.NewRequest("GET", "/urlshort", nil)
	response := httptest.NewRecorder()

	handler, err := JSONHandler([]byte(json), fallback)

	assert.Nil(t, err)

	handler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "https://github.com/gophercises/urlshort", response.HeaderMap.Get("Location"))
}
