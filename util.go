package main

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	_ "github.com/skymkmk/ShortLink/statik"
)

var tlds []string
var domain_tld_regexp *regexp.Regexp

const DOMAIN_TLD_REGEX = `\w+$`

type NewShortLinkResp struct {
	Status  int    `json:"status"`
	Err     string `json:"error"`
	RealURL string `json:"realURL"`
}

func isOverflow(num byte) bool {
	return num+1 == 0
}

func initTlds() {
	resp, err := http.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, _ := io.ReadAll(resp.Body)
	for _, line := range strings.Split(string(data), "\n") {
		tld := strings.TrimSpace(line)
		if tld != "" && !strings.HasPrefix(tld, "#") {
			tlds = append(tlds, tld)
		}
	}
	domain_tld_regexp, _ = regexp.Compile(DOMAIN_TLD_REGEX)
}

func checkDomainAvailablity(domain string) bool {
	if len(tlds) == 0 {
		log.Fatal("No TLDs")
	}
	tld := domain_tld_regexp.FindString(domain)
	for _, value := range tlds {
		if strings.EqualFold(tld, value) {
			return true
		}
	}
	return false
}

func URLInvalid(c *gin.Context) {
	c.JSON(400, NewShortLinkResp{
		Status:  1,
		Err:     "Invalid URL",
		RealURL: "",
	})
}

func serveWeb(c *gin.Context) {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	http.FileServer(statikFS).ServeHTTP(c.Writer, c.Request)
}
