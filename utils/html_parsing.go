package utils

import (
	"io"
	"net/http"
	"regexp"
)

var originUrlPattern = regexp.MustCompile(`(https?://[^/]*)/`)
var aHrefLinkPattern = regexp.MustCompile(`<a *href *= *['"](?P<link>[^#'"]*)['"]`)

// dunno why there is proxy by default
var defaultTransport http.RoundTripper = &http.Transport{Proxy: nil}
var client = &http.Client{Transport: defaultTransport}

func ParseOriginUrl(url string) string {
	matches := originUrlPattern.FindAllStringSubmatch(url, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	} else {
		return ""
	}
}

func ParseAllLinks(html string) []string {
	matches := aHrefLinkPattern.FindAllStringSubmatch(html, -1)
	var allAddrs []string
	if len(matches) > 0 {
		for _, groups := range matches {
			if len(groups) > 1 {
				allAddrs = append(allAddrs, groups[1])
			}
		}
		return allAddrs
	} else {
		return []string{}
	}
}

func GetHtmlByUrl(url string) string {
	resp, e1 := client.Get(url)
	if e1 != nil {
		println(e1.Error())
		return ""
	} else {
		defer resp.Body.Close()
		body, e2 := io.ReadAll(resp.Body)
		if e2 != nil {
			println(e2.Error())
			return ""
		} else {
			return string(body)
		}
	}
}
