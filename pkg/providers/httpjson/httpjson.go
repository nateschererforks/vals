package httpjson

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xpath"

	"github.com/nateschererforks/vals/pkg/api"
	"github.com/nateschererforks/vals/pkg/log"
)

type provider struct {
	// Keeping track of httpjson services since we need a service per url
	log        *log.Logger
	docs       map[string]*jsonquery.Node
	protocol   string
	floatAsInt bool
}

func New(l *log.Logger, cfg api.StaticConfig) *provider {
	p := &provider{
		log: l,
	}

	// Should the protocol be insecure i.e. http
	insecureArg := cfg.String("insecure")
	p.protocol = "https"
	if insecureArg == "true" {
		p.protocol = "http"
	}

	// By default JSON will return large integers as float64
	floatAsIntArg := cfg.String("floatAsInt")
	p.floatAsInt = false
	if floatAsIntArg == "true" {
		p.floatAsInt = true
	}

	// Initialize docs map to store the json object for use multiple times
	if len(p.docs) == 0 {
		p.docs = make(map[string]*jsonquery.Node)
	}

	return p
}

func GetXpathFromUri(uri string) (xpathExpression string, err error) {
	paths := strings.Split(uri, "#/")
	if len(paths) == 1 {
		return "", fmt.Errorf("no xpath expression found in uri: %s", uri)
	}
	_, err = xpath.Compile(paths[1])
	if err != nil {
		return "", fmt.Errorf("unable to compile xpath expression '%s' from uri: %s", xpathExpression, uri)
	}
	xpathExpression = paths[1]

	return xpathExpression, nil
}

func GetUrlFromUri(uri string, protocol string) (string, error) {
	// Remove httpjson:// prefix
	trimmedUrl := strings.TrimPrefix(uri, "httpjson://")
	// Remove xpath suffix
	trimmedUrl = strings.Split(trimmedUrl, "#/")[0]

	fullURL := fmt.Sprintf("%s://%s", protocol, trimmedUrl)

	parsedUrl, parseErr := url.Parse(fullURL)
	if parseErr != nil {
		return "", fmt.Errorf("invalid URL: %s, error: %s", fullURL, parseErr.Error())
	}

	if parsedUrl.Host == "" {
		return "", fmt.Errorf("no domain found in uri: %s", uri)
	}

	query := parsedUrl.Query()

	for _, key := range []string{"insecure", "floatAsInt"} {
		query.Del(key)
	}

	parsedUrl.RawQuery = query.Encode()

	return parsedUrl.String(), nil
}

func (p *provider) GetJsonDoc(url string) error {
	if _, ok := p.docs[url]; !ok {
		doc, err := jsonquery.LoadURL(url)
		if err != nil {
			return fmt.Errorf("error fetching json document at %v: %v", url, err)
		}
		p.log.Debugf("httpjson: successfully retrieved JSON data from: %s", url)
		p.docs[url] = doc
	}

	return nil
}

func (p *provider) GetString(uri string) (string, error) {
	url, err := GetUrlFromUri(uri, p.protocol)
	if err != nil {
		return "", err
	}
	err = p.GetJsonDoc(url)
	if err != nil {
		return "", err
	}
	xpathQuery, err := GetXpathFromUri(uri)
	if err != nil {
		return "", err
	}

	returnValue := ""
	var values []string
	node, err := jsonquery.Query(p.docs[url], xpathQuery)
	if err != nil {
		return "", fmt.Errorf("error querying the xpath expression %s, error: %s", xpathQuery, err)
	}

	if node == nil {
		return "", fmt.Errorf("unable to query doc for value with xpath query using %v", uri)
	}

	if node.FirstChild.Data != node.LastChild.Data {
		return "", fmt.Errorf("location %v has child nodes at %v, please use a more granular query", xpathQuery, url)
	}

	childNodesLength := countChildNodes(node)

	if childNodesLength > 1 {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			values = append(values, child.Value().(string))
		}
		returnValue = strings.Join(values, ",")
	} else {
		returnValue = node.FirstChild.Value().(string)
	}

	if p.floatAsInt {
		intValue, err := strconv.ParseFloat(returnValue, 64)
		if err != nil {
			return "", fmt.Errorf("unable to convert possible float to int for value: %v", returnValue)
		}
		returnValue = fmt.Sprintf("%.0f", intValue)
	}

	return returnValue, nil
}

func countChildNodes(node *jsonquery.Node) int {
	// Check if there are more child nodes i.e. keys under this json key
	count := 0
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		count++
	}
	return count
}

func (p *provider) GetStringMap(key string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("we should not be in the GetStringMap method")
}
