package certsrv

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

func ParseHTMLResponse(response io.Reader) (string, error) {
	tokenizer := html.NewTokenizer(response)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			fmt.Println("End")
			break
		}
		tag, hasAttr := tokenizer.TagName()
		if !hasAttr {
			continue
		}
		if string(tag) != "a" {
			continue
		}
		for {
			attrKey, attrValue, moreAttr := tokenizer.TagAttr()
			if string(attrKey) == "href" {
				href := string(attrValue)
				if !strings.Contains(href, "b64") {
					continue
				}
				if strings.Contains(href, "p7b") || strings.Contains(href, "CA") {
					continue
				}
				return href, nil
			}
			if !moreAttr {
				break
			}
		}
	}
	return "", errors.New("No valid link found")
}
