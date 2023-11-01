package url

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// URLInfo contains information about a URL.
type URLInfo struct {
	Scheme   string `json:"scheme"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	Query    string `json:"query"`
	Fragment string `json:"fragment"`
}

// Sitemap contains the generated sitemap.
type Sitemap struct {
	Schemes map[string]map[string]map[string]map[string][]string `json:"schemes"`
}

// Init initializes an empty sitemap.
func Init() *Sitemap {
	return &Sitemap{
		Schemes: make(map[string]map[string]map[string]map[string][]string),
	}
}

// ValidateURL validates whether a given URL is valid.
func ValidateURL(URL string, filterParam bool) bool {
	regex := `^[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`
	match, err := regexp.MatchString(regex, URL)
	if filterParam && err == nil && match {
		return regexp.MustCompile(`\?[A-Za-z0-9]+=`).MatchString(URL)
	}
	return err == nil && match
}

// GetFromFile reads a text file and returns the found URLs.
func GetFromFile(filepath string) ([]string, error) {
	var URLs []string
	file, err := os.Open(filepath)
	if err != nil {
		return URLs, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		URL := scanner.Text()
		URLs = append(URLs, URL)
	}

	if err := scanner.Err(); err != nil {
		return URLs, err
	}

	return URLs, nil
}

// SaveToJson saves the sitemap in JSON format to the specified file.
func (s *Sitemap) SaveToJson(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(s); err != nil {
		return err
	}

	return nil
}

// Add adds the URLs to the sitemap and validates them.
func (s *Sitemap) Add(URLs []string) error {
	for _, URL := range URLs {
		if ValidateURL(URL, false) {
			s.addURL(URL)
		} else {
			fmt.Printf("Invalid URL: %s\n", URL)
		}
	}
	return nil
}

// Show displays the complete sitemap line by line.
func (s *Sitemap) Show() {
	for scheme, domains := range s.Schemes {
		for domain, paths := range domains {
			for path, queries := range paths {
				for _, query := range queries["queries"] {
					fmt.Println(scheme + "://" + domain + path + "?" + query)
				}
			}
		}
	}
}

// Show displays the sitemap hierarchically.
func (s *Sitemap) ShowTree() {
	for scheme, domains := range s.Schemes {
		fmt.Printf("Scheme: %s\n", scheme)
		for domain, paths := range domains {
			fmt.Printf("  Domain: %s\n", domain)
			for path, queries := range paths {
				fmt.Printf("    Path: %s\n", path)
				fmt.Printf("      Queries: %s\n", strings.Join(queries["queries"], ", "))
			}
		}
	}
}

// addURL adds a valid URL to the sitemap.
func (s *Sitemap) addURL(URL string) {
	parsedURL, _ := url.Parse(URL)
	urlInfo := URLInfo{
		Scheme:   parsedURL.Scheme,
		Domain:   parsedURL.Host,
		Path:     parsedURL.Path,
		Query:    parsedURL.RawQuery,
		Fragment: parsedURL.Fragment,
	}

	scheme := s.Schemes[urlInfo.Scheme]
	if scheme == nil {
		scheme = make(map[string]map[string]map[string][]string)
		s.Schemes[urlInfo.Scheme] = scheme
	}

	domain := scheme[urlInfo.Domain]
	if domain == nil {
		domain = make(map[string]map[string][]string)
		scheme[urlInfo.Domain] = domain
	}

	path := domain[urlInfo.Path]
	if path == nil {
		path = make(map[string][]string)
		domain[urlInfo.Path] = path
	}

	queries := path["queries"]
	queries = append(queries, urlInfo.Query)
	path["queries"] = queries
}

// AddValueParam adds the specified value to the corresponding queries in the sitemap.
func (s *Sitemap) AddValueParam(key, value string) {
	for _, schemeData := range s.Schemes {
		for _, domainData := range schemeData {
			for _, pathData := range domainData {
				for i, query := range pathData["queries"] {
					queries := []string{}
					pairs := strings.Split(query, "&")
					for _, pair := range pairs {
						parts := strings.SplitN(pair, "=", 2)
						if len(parts) == 2 && (parts[0] == key || key == "*") {
							queries = append(queries, parts[0]+"="+parts[1]+value)
						} else if len(parts) == 1 && (pair == key || key == "*") {
							queries = append(queries, pair+"="+value)
						} else {
							queries = append(queries, pair)
						}
					}
					pathData["queries"][i] = strings.Join(queries, "&")
				}
			}
		}
	}
}

// SetValueParam sets the specified value to the corresponding queries in the sitemap.
func (s *Sitemap) SetValueParam(key, value string) {
	for _, schemeData := range s.Schemes {
		for _, domainData := range schemeData {
			for _, pathData := range domainData {
				for i, query := range pathData["queries"] {
					queries := []string{}
					pairs := strings.Split(query, "&")
					for _, pair := range pairs {
						parts := strings.SplitN(pair, "=", 2)
						if len(parts) == 2 && (parts[0] == key || key == "*") {
							queries = append(queries, parts[0]+"="+value)
						} else if len(parts) == 1 && (pair == key || key == "*") {
							queries = append(queries, pair+"="+value)
						} else {
							queries = append(queries, pair)
						}
					}
					pathData["queries"][i] = strings.Join(queries, "&")
				}
			}
		}
	}
}

// MergeKeysParam merges queries of the same path into a single query.
func (s *Sitemap) MergeKeysParam() {
	for _, schemeData := range s.Schemes {
		for _, domainData := range schemeData {
			for _, pathData := range domainData {
				mergedQueries := make(map[string][]string)
				for _, query := range pathData["queries"] {
					queryParts := strings.Split(query, "&")
					for _, queryPart := range queryParts {
						kv := strings.Split(queryPart, "=")
						if len(kv) != 2 {
							continue
						}
						key, value := kv[0], kv[1]
						mergedQueries[key] = append(mergedQueries[key], value)
					}
				}
				// Recreate a single query with all merged parameters
				var newQueries []string
				for key, values := range mergedQueries {
					newQuery := key + "=" + values[0]
					newQueries = append(newQueries, newQuery)
				}
				pathData["queries"] = []string{strings.Join(newQueries, "&")}
			}
		}
	}
}
