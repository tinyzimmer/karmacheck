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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
)

func checkSubs(arg *string) (subs []string, err error) {

	// Checks that subreddit is defined. More checks can be added here later.
	subs = strings.Split(*arg, ",")
	if len(subs) == 0 {
		err = errors.New(NO_SUBREDDIT_ERROR)
	}
	return

}

func kdIsConfident(content []byte) (res bool) {

	re := regexp.MustCompile(MARKDOWN_LINK_REGEX)
	data := re.FindAll(content, -1)
	if len(data) < KARMA_DECAY_COMMENT_LINKS_CONFIDENCE {
		res = false
	} else {
		res = true
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

func checkKarmaDecay(p reddit.Post) (resp string, err error) {

	// Get the KarmaDecay URL

	karmaUrl := fmt.Sprintf("%s%s", KARMA_DECAY_URL, p.Permalink)
	if err != nil {
		return
	}
	log.Printf("Checking KarmaDecay for: %s\n", p.Permalink)
	log.Printf("Author: u/%s\n", p.Author)
	log.Printf("Title: %s\n", p.Title)
	log.Printf("Created: %s", time.Unix(int64(p.CreatedUTC), 0))

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
