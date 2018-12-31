package urlshort

import (
	"net/http"

	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		path := request.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(response, request, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(response, request)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return parseAndBuildMapHandler(parseYAML, yamlBytes, fallback)
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//     [{
//         "path": "/some-path",
//         "url": "https://www.some-url.com/demo"
//     }]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return parseAndBuildMapHandler(parseJSON, jsonBytes, fallback)
}

// Represents the mapping between a path and the target url, as represented in JSON/YAML.
type pathMapping struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// An alias for a function that parses mappings from a byte sequence (it can be from YAML, JSON, etc)
type dataParser func([]byte) ([]pathMapping, error)

// A dataParser to parse the mappings from a YAML.
func parseYAML(yamlBytes []byte) ([]pathMapping, error) {
	var mappings []pathMapping
	if err := yaml.Unmarshal(yamlBytes, &mappings); err != nil {
		return nil, err
	}

	return mappings, nil
}

// A dataParser to parse the mappings from a JSON.
func parseJSON(jsonBytes []byte) ([]pathMapping, error) {
	var pathMappings []pathMapping
	if err := json.Unmarshal(jsonBytes, &pathMappings); err != nil {
		return nil, err
	}

	return pathMappings, nil
}

// Method to transform a slice of mappings to a map indexed by the path, to perform faster lookups.
func buildMap(pathMappings []pathMapping) map[string]string {
	pathsToUrls := make(map[string]string)

	for _, mapping := range pathMappings {
		pathsToUrls[mapping.Path] = mapping.URL
	}

	return pathsToUrls
}

// Contains the common logic for handlers that parse a byte sequence. It receives the dataParser method
// so it can parse mappings, transform them into a map and return a MapHandler for this data.
func parseAndBuildMapHandler(parser dataParser, data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathMappings, err := parser(data)
	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMap(pathMappings)
	return MapHandler(pathsToUrls, fallback), nil
}
