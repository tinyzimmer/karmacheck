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
	RSS_URL_FORMAT                = "%s/r/%s/new/%s"
	REQUEST_AGENT_HEADER          = "User-Agent"
	REQUEST_AGENT                 = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
	NO_CONTENT_ERROR              = "KarmaDecay could not locate any media in the post"
	NO_SIMILAR_POSTS_ERROR        = "KarmaDecay could not find any similar posts"
	MALFORMED_URL_ERROR           = "Malformed URL: %s"
	NO_SUBREDDIT_ERROR            = "No subreddit specified"
	KARMA_DECAY_NO_CONTENT_STRING = "Unable to find an image"
	LOCAL_FOUND_MATCHES           = "Found matches. Below is the reddit comment text."
	MARKDOWN_SEARCH_REGEX         = "Anyone[^<]*"
	MARKDOWN_VALID_CHECK          = "[Source: karmadecay]"
)

var (
	subreddit = flag.String("s", "", "Subreddit to watch")
)

// Type definitions for Reddit RSS feed parsing

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
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
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

	// Checks that subreddit is defined. More checks can be added here later.

	if *sub == "" {
		err = errors.New(NO_SUBREDDIT_ERROR)
	}
	return
}

func getMarkdownComment(content []byte) (comment string) {

	// Uses the regex defined above to fetch the markdown comment from the
	// KarmaDecay page.

	re := regexp.MustCompile(MARKDOWN_SEARCH_REGEX)
	data := re.Find(content)
	comment = string(data)
	if !strings.Contains(comment, MARKDOWN_VALID_CHECK) { // title flipped the regex
		comment = ""
	}
	return

}

func hasContent(content []byte) (res bool) {

	// Checks KarmaDecay response for presence of the error specifying
	// no media content could be found in the post

	if strings.Contains(string(content), KARMA_DECAY_NO_CONTENT_STRING) {
		res = false
	} else {
		res = true
	}
	return

}

func getCommentUrl(fullUrl string) (commentUrl string) {

	// removes https and the domain from a full URL. Just showing the
	// subreddit/comment part of it.

	splitUrl := strings.Split(fullUrl, "/")
	commentUrl = strings.Join(splitUrl[3:], "/")
	return

}

func getUrl(url string) (response []byte, err error) {

	// Does a GET request against a URL using the User-Agent header defined above.
	// Returns the raw bytes of the respose

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

func checkKarmaDecay(entry Entry) (resp string, err error) {

	// Make sure our URL is legit

	if !strings.Contains(entry.Link.Href, "reddit") {
		err = errors.New(fmt.Sprintf(MALFORMED_URL_ERROR, entry.Link.Href))
		return
	}

	// Replace the reddit part of the URL with karmadecay

	karmaUrl := strings.Replace(entry.Link.Href, REDDIT_URL, KARMA_DECAY_URL, 1)
	commentUrl := getCommentUrl(karmaUrl)
	log.Printf("Checking KarmaDecay for: %s\n", commentUrl)
	log.Printf("Author: %s\n", entry.Author.Name)
	log.Printf("Title: %s\n", entry.Title)
	log.Printf("Link: %s", entry.Link.Href)
	log.Printf("Created: %s", entry.Updated)

	// Get the content of the KarmaDecay page for that post. May take 10-20
	// seconds on a post that is indeed OC
	response, err := getUrl(karmaUrl)
	if err != nil {
		return
	}

	// Ensure KarmaDecay detected content in the page, and fetch the markdown

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

	// Poll the reddit RSS feed for the latest submissions from a subreddit

	url := fmt.Sprintf(RSS_URL_FORMAT, REDDIT_URL, *subreddit, RSS_ARG)
	resp, err := getUrl(url)
	xml.Unmarshal(resp, &feed)
	return

}

func main() {

	flag.Parse()
	err := checkSub(subreddit) // make sure subreddit is defined
	if err != nil {
		log.Fatal(err)
	}
	data, err := getLatestSubmissions(subreddit) // get latest reddit posts
	if err != nil {
		log.Fatal(err)
	}
	// check the most recent against KarmaDecay
	res, err := checkKarmaDecay(data.Entries[0])
	if err != nil {
		log.Println(err)
	} else {
		log.Println(LOCAL_FOUND_MATCHES)
		fmt.Println(res)
	}

}
