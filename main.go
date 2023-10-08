//go:generate bash -c "cd html && npm install && npm run build"
//go:generate statik -src=./html/dist/

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58"
)

func getShortLink(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		serveWeb(c)
	}
	realURL, err := queryCode(shortCode)
	if err != nil {
		serveWeb(c)
	} else {
		u, _ := url.Parse(realURL)
		if u.Scheme == "" {
			c.Redirect(301, "http://"+realURL)
		} else {
			c.Redirect(301, realURL)
		}
	}
}

func newShortLink(c *gin.Context) {
	realUrl, err := url.QueryUnescape(c.Query("url"))
	if err != nil {
		URLInvalid(c)
		return
	}
	u, err := url.Parse(realUrl)
	fmt.Println(u.Host)
	if err != nil {
		URLInvalid(c)
		return
	} else if u.Host == "" {
		u, _ = url.Parse("http://" + realUrl)
		if u.Host == "" {
			URLInvalid(c)
			return
		}
	}
	if !checkDomainAvailablity(u.Host) {
		URLInvalid(c)
		return
	}
	url2bytes := []byte(realUrl)
	for {
		hash := sha256.Sum256(url2bytes)
		hexHash := hex.EncodeToString(hash[:])
		encoded := base58.Encode([]byte(hexHash))[0:6]
		urlInDB, err := queryCode(encoded)
		if err != nil {
			if urlInDB == "" {
				c.Writer.WriteString(err.Error())
				return
			} else {
				err = insertURL(encoded, realUrl)
				if err != nil {
					c.JSON(500, NewShortLinkResp{
						Status:  2,
						Err:     err.Error(),
						RealURL: "",
					})
					return
				}
				c.JSON(200, NewShortLinkResp{
					Status:  0,
					Err:     "",
					RealURL: c.Request.Host + "/" + encoded,
				})
				break
			}
		} else {
			if realUrl == urlInDB {
				c.JSON(200, NewShortLinkResp{
					Status:  3,
					Err:     "Already added!",
					RealURL: c.Request.Host + "/" + encoded,
				})
				return
			} else {
				p := len(url2bytes) - 1
				for {
					if isOverflow(url2bytes[p]) {
						url2bytes[p] = 0
						if p == 0 {
							break
						} else {
							p--
						}
					} else {
						url2bytes[p] += 1
						break
					}
				}
			}
		}
	}
}

func main() {
	initTlds()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	initSQL()
	router.Handle("GET", "/", getShortLink)
	router.Handle("GET", "/:shortCode", getShortLink)
	router.Handle("GET", "/:shortCode/*any", getShortLink)
	router.Handle("GET", "/api/v1/newShortLink", newShortLink)
	router.Run(port())
}
