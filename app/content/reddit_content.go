/*
 *   Copyright 2016 Rémy MATHIEU
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package content

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RedditContent struct {
	Url string
}

func (reddit RedditContent) callAPI() ([]Content, error) {
	resp, err := getWebserviceResponse(reddit.Url)
	if err != nil {
		log.Println("Error while calling url: " + reddit.Url)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Println("Error while calling API")
		return nil, err
	}

	rv := make([]Content, 0)

	doc.Find("div.Post").Each(func(i int, selec *goquery.Selection) {

		// ignore sticky posts
		if selec.HasClass("stickied") {
			return
		}

		if len(rv) > 20 {
			return
		}

		t := selec.Find("h3")
		title := t.First().Find("span").First().Text()

		// Limit size of content
		if len(title)+urlLength > 280 {
			title = title[0:279-urlLength] + "…"
		}

		l := selec.Find("a[target=\"_blank\"]")
		externalLink, _ := l.First().Attr("href")

		// self posts
		if strings.HasPrefix(externalLink, "/r/") {
			externalLink = "https://reddit.com" + externalLink
		}

		rv = append(rv, Content{
			Text: title,
			Url:  externalLink,
		})
	})

	return rv, nil
}
