package urlshort

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		fmt.Printf("incoming request url: %s\n", url)

		// check if path is in pathsToURls mapper
		val, exists := pathsToUrls[url]
		if exists {
			http.Redirect(w, r, val, http.StatusSeeOther)
			return
		}
		fmt.Println("couldn't map path, using fallback handler")
		fallback.ServeHTTP(w, r)
		return
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
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	type PathUrlMap map[string]string
	var urlMaps []PathUrlMap

	// parse yaml into a slice of path to url mappings
	err := yaml.Unmarshal(yml, &urlMaps)
	if err != nil {
		err = fmt.Errorf("Error unmarshalling YAML: %w\n", err)
		return nil, err
	}
	pathsToUrls := make(map[string]string)

	// loop over urlMaps and condense into one map
	for _, m := range urlMaps {
		path, _ := m["path"]
		url, _ := m["url"]
		pathsToUrls[path] = url
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fallback.ServeHTTP(w, r)
	}, nil
}
