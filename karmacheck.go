/**
    This file is part of KarmaCheck.

    KarmaCheck is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    KarmaCheck is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with KarmaCheck.  If not, see <http://www.gnu.org/licenses/>.
**/

package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	REDDIT_URL                    = "https://www.reddit.com"
	KARMA_DECAY_URL               = "https://www.karmadecay.com"
	RSS_ARG                       = ".rss"
	REQUEST_AGENT_HEADER          = "User-Agent"
	REQUEST_AGENT                 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
	NO_CONTENT_ERROR              = "KarmaDecay could not locate any media in the post"
	NO_SIMILAR_POSTS_ERROR        = "KarmaDecay could not find any similar posts"
	MALFORMED_URL_ERROR           = "Malformed URL: %s"
	NO_SUBREDDIT_ERROR            = "No subreddit specified"
	KARMA_DECAY_NO_CONTENT_STRING = "Unable to find an image"
	LOCAL_FOUND_MATCHES           = "Found matches. Below is the reddit comment text."
	MARKDOWN_SEARCH_REGEX         = "Anyone[^<]*"
)

var (
	subreddit = flag.String("s", "", "Subreddit to watch")
)

type Feed struct {
	Namespace string   `xml:"xmlns,attr"`
	Category  Category `xml:"category"`
	Updated   string   `xml:"updated"`
	Icon      string   `xml:"icon"`
	Id        string   `xml:"id"`
	Links     []Link   `xml:"link"`
	Title     string   `xml:"title"`
	Entries   []Entry  `xml:"entry"`
}

type Category struct {
	Term  string `xml:"term,attr"`
	Label string `xml:"label,attr"`
}

type Link struct {
	Rel      string   `xml:"rel,attr"`
	Href     string   `xml:"href,attr"`
	Type     string   `xml:"type,attr"`
	Category Category `xml:"category"`
}

type Entry struct {
	Author  Author `xml:"author"`
	Content string `xml:"content"`
	Id      string `xml:"id"`
	Link    Link   `xml:"link"`
	Updated string `xml:"updated"`
	Title   string `xml:"title"`
}

type Author struct {
	Name string `xml:"name"`
	Uri  string `xml:"uri"`
}

func checkSub(sub *string) (err error) {
	if *sub == "" {
		err = errors.New(NO_SUBREDDIT_ERROR)
	}
	return
}

func getMarkdownComment(content []byte) (comment string) {

	re := regexp.MustCompile(MARKDOWN_SEARCH_REGEX)
	data := re.Find(content)
	comment = string(data)
	return

}

func hasContent(content []byte) (res bool) {

	if strings.Contains(string(content), KARMA_DECAY_NO_CONTENT_STRING) {
		res = false
	} else {
		res = true
	}
	return

}

func getCommentUrl(fullUrl string) (commentUrl string) {
	splitUrl := strings.Split(fullUrl, "/")
	commentUrl = strings.Join(splitUrl[3:], "/")
	return
}

func getUrl(url string) (response []byte, err error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set(REQUEST_AGENT_HEADER, REQUEST_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	return

}

func checkKarmaDecay(redditFullUrl string) (resp string, err error) {

	if !strings.Contains(redditFullUrl, "reddit") {
		err = errors.New(fmt.Sprintf(MALFORMED_URL_ERROR, redditFullUrl))
		return
	}
	karmaUrl := strings.Replace(redditFullUrl, REDDIT_URL, KARMA_DECAY_URL, 1)
	commentUrl := getCommentUrl(karmaUrl)
	log.Printf("Checking KarmaDecay for: %s\n", commentUrl)
	response, err := getUrl(karmaUrl)
	if err != nil {
		return
	}
	if !hasContent(response) {
		err = errors.New(NO_CONTENT_ERROR)
	} else {
		resp = getMarkdownComment(response)
		if len(resp) == 0 {
			err = errors.New(NO_SIMILAR_POSTS_ERROR)
		}
	}
	return

}

func getLatestSubmissions(subreddit *string) (feed *Feed, err error) {

	url := fmt.Sprintf("%s/r/%s/new/%s", REDDIT_URL, *subreddit, RSS_ARG)
	resp, err := getUrl(url)
	xml.Unmarshal(resp, &feed)
	return

}

func main() {

	flag.Parse()
	err := checkSub(subreddit)
	if err != nil {
		log.Fatal(err)
	}
	data, err := getLatestSubmissions(subreddit)
	if err != nil {
		log.Fatal(err)
	}
	res, err := checkKarmaDecay(data.Entries[0].Link.Href)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(LOCAL_FOUND_MATCHES)
		fmt.Println(res)
	}

}
