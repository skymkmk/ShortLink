package main

import (
	"net/url"
	"regexp"
)

func trackRemover(link string) string {
	var finalURL string
	var err error
	link2byte := []byte(link)
	isb23, _ := regexp.Match(".*b23\\.tv.*", link2byte)
	isPilipili, _ := regexp.Match(".*bilibili\\.com.*", link2byte)
	if isb23 {
		finalURL, err = b23dottv(link)
		if err == nil {
			return finalURL
		}
	}
	if isPilipili {
		finalURL = pilipili(link)
		return finalURL
	}
	return link
}

func b23dottv(link string) (string, error) {
	finalURL, err := getFinalURL(link)
	if err != nil {
		return "", err
	}
	pureURL := pilipili(finalURL)
	return pureURL, err
}

func pilipili(link string) string {
	return removeAllParameter(link)
}

func removeAllParameter(link string) string {
	url, _ := url.Parse(link)
	return "https://" + url.Host + url.Path
}
