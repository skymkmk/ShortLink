package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
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

func port() string {
	var port int = 3000
	var err error
	arguments := os.Args[1:]
	for idx, val := range arguments {
		if val == "-p" || val == "-port" {
			if idx != len(arguments)-1 {
				port, err = strconv.Atoi(arguments[idx+1])
				if err != nil {
					port = 3000
				}
			}
			break
		}
	}
	return ":" + strconv.Itoa(port)
}

func isURL(theURL string) bool {
	u, err := url.Parse(theURL)
	if err != nil {
		return false
	}
	if u.Host == "" {
		if isMagnetURL(theURL) {
			return true
		}
		u, err = url.Parse("http://" + theURL)
		if err != nil || u.Host == "" {
			return false
		}
	}
	if !checkDomainAvailablity(u.Host) {
		return false
	}
	return true
}

func isMagnetURL(theURL string) bool {
	u, err := url.Parse(theURL)
	if err != nil {
		return false
	}
	if u.Scheme == "magnet" {
		return u.Query().Has("xt")
	}
	return false
}
